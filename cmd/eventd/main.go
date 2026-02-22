package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/events"
	"starbase.ag/liquidity/liquid/utils"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/pkg/sentry"
)

func main() {
	configFn := "./config.toml"

	app := &cli.App{
		Name:  "eventd",
		Usage: "Sync blockchain event to Kafka",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "conf",
				Value:       "./config.toml",
				Usage:       "config file",
				Destination: &configFn,
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err) // err)
	}
}

func run(ctx *cli.Context) error {
	cfg, err := config.LoadConfigAWS(ctx.String("conf"), config.CONFIG_AWS_ENV_EVENTS)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	logger.Init(cfg.SentryOpts.DSN, cfg.SentryOpts.Env, cfg.IsLocal(), &cfg.Log)

	subscriber, err := events.NewEventSubscirber(
		cfg.Chain.URL,
		cfg.Chain.WSS,
		cfg.Chain.Subscribe,
		uint64(cfg.Chain.PollingInterval),
		cfg.Chain.BlockInterval,
		0,
		uint32(cfg.Chain.PollingSteps)) // nolint
	if err != nil {
		panic(err.Error())
	}

	eventsCh := make(chan []types.Log, 100)

	startBlock, lastEvt, err := events.FetchLatestBlock(cfg.Kafka.Topic, cfg.Kafka.Brokers, "producer-resumer") // uint64(0)
	if err != nil {
		panic(err.Error())
	}

	go func() {
		if startBlock > 0 {
			logger.Info().Msgf("resume events from block %d logIndex=%d", startBlock, lastEvt.Index)
			// fillup events
			evts, err := subscriber.GetLogsFromToReturn(context.Background(), nil, startBlock, startBlock, 1)
			if err != nil {
				utils.PanicWithSentry("GetLogsFromToReturn failed: block=%d error=%v", startBlock, err)
			}

			filled := 0

			if len(evts) > 0 {
				if evts[len(evts)-1].Index == lastEvt.Index {
					// whole block events has been published
				} else {
					index := 0

					for i, evt := range evts {
						if evt.Index > lastEvt.Index {
							index = i
							break
						}
					}

					filled = len(evts) - index
					eventsCh <- evts[index:]
				}
			}

			logger.Info().Msgf("fillup events: %d", filled)

			startBlock++
		}

		err := subscriber.SubscribeEvents(context.Background(), startBlock, []string{}, eventsCh)
		if err != nil {
			panic("subscribe events failed: " + err.Error())
		}
	}()

	if cfg.SentryOpts.DSN != "" {
		_ = sentry.InitSentry(cfg.SentryOpts.DSN, cfg.SentryOpts.Env, cfg.SentryOpts.Debug)
	}

	return events.ProduceEvents(context.Background(), cfg.Kafka.Topic, cfg.Kafka.Brokers, eventsCh)
}
