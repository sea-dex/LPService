package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/events"
	"starbase.ag/liquidity/liquid/handlers"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/liquid/utils"
	"starbase.ag/liquidity/pkg/logger"
)

func main() {
	configFn := "./config.toml"
	start := int64(-1)
	mode := "kafka" // rpc or kafka

	app := &cli.App{
		Name:  "liquidity",
		Usage: "Consumer kafka events, parse events, track liquidity",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "conf",
				Value:       "./config.toml",
				Usage:       "config file",
				Destination: &configFn,
			},
			&cli.Int64Flag{
				Name:        "start",
				Value:       -1,
				Usage:       "parse from kafka offset(Not block number)",
				Destination: &start,
			},
			&cli.StringFlag{
				Name:        "mode",
				Value:       "kafka", // kafka or rpc
				Usage:       "get events from RPC or kafka",
				Destination: &mode,
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

	// 1371680 create uniswap v3 factory
	// 2284119 create first pool 0x9c625cb2ae462515fc614528f727d1a4d3bfbde2
	// 2315820
	// startBlock := uint64(2315820) // uniswap v3 factory
	// go func() {
	// 	err := subscriber.SubscribeEvents(context.Background(), startBlock, []string{}, eventsCh)
	// 	if err != nil {
	// 		panic("subscribe events failed: " + err.Error())
	// 	}
	// }()
	//
	// subscriber.GetLogsFromTo(context.Background(), nil, 2284119, 2284119, 100, eventCh)

	var wg sync.WaitGroup

	parser := buildParser()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	cctx, cancel := context.WithCancel(ctx.Context)

	wg.Add(1)

	mode := strings.ToLower(ctx.String("mode"))
	println("mode:", mode)

	go func() {
		logger.Info().Msgf("start liquidity mode: %s", mode)

		if mode == "kafka" {
			fn := func(events []types.Log) error {
				parser.ParseEvents(events)
				return nil
			}

			err = events.ConsumeEvents(cctx, cfg.Kafka.Topic, cfg.Kafka.Brokers, cfg.Kafka.Group, fn)
			if err != nil {
				utils.PanicWithSentry("start consumer failed: %v", err)
			}
		} else if mode == "rpc" {
			// 1371680 create uniswap v3 factory
			// 2112314 create pool=0x4c36388be6f416a29c8d8eee81c771ce6be14b18
			// 2284119 create pool 0x9c625cb2ae462515fc614528f727d1a4d3bfbde2
			// 2315820
			startBlock := uint64(2112314)
			eventsCh := make(chan []types.Log, 100)

			go func() {
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

				err = subscriber.SubscribeEvents(context.Background(), startBlock, []string{}, eventsCh)
				if err != nil {
					panic("subscribe events failed: " + err.Error())
				}
			}()

			// handle events
			go func() {
				for events := range eventsCh {
					if !cfg.IsProd() {
						logger.Info().Uint64("startBlock", events[0].BlockNumber).Msg("parse events")
					}

					parser.ParseEvents(events)
				}
			}()
		} else {
			panic(fmt.Sprintf("invalid param mode: %s, either rpc or kafka", mode))
		}

		wg.Done()
	}()

	<-c
	cancel()
	logger.Info().Msg("recv signal, terminating liquidity ....")
	wg.Wait()
	logger.Info().Msg("liquidity exit")

	return nil
}

func buildParser() *handlers.Parser {
	pool.InitABIs()

	parser := handlers.NewParser()
	uniswapv3 := handlers.CreateCAMMHandler("UniswapV3", common.UniswapV3FactoryAddress)
	parser.AddLiquidor(uniswapv3)

	return parser
}
