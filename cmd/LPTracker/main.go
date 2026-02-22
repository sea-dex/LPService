package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" //nolint
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/liquid/events"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/liquid/swapor"
	"starbase.ag/liquidity/pkg/health"
	"starbase.ag/liquidity/pkg/logger"
	"starbase.ag/liquidity/version"
)

func main() {
	version.PrintVersion(os.Stdout)

	configFn := "./config.toml"

	app := &cli.App{
		Name:  "lpservice",
		Usage: "Consumer events, parse events, track pool liquidity",
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

	go func() {
		server := &http.Server{
			Addr:              ":6060",
			ReadHeaderTimeout: 10 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			panic("start pprof server failed: " + err.Error())
		}
	}()

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

	pool.InitABIs()

	subscriber, err := events.NewEventSubscirber(
		cfg.Chain.URL,
		cfg.Chain.WSS,
		cfg.Chain.Subscribe,
		uint64(cfg.Chain.PollingInterval),
		cfg.Chain.BlockInterval,
		0,
		uint32(cfg.Chain.PollingSteps)) // nolint
	if err != nil {
		logger.Fatal().Msg("create subscriber failed: " + err.Error())
	}

	// create event consumer
	es := swapor.CreateEventHandler(&cfg)
	// start subscribe
	cctx, cancel := context.WithCancel(ctx.Context)
	go subscriber.SubscribeBlock(cctx)
	go func() {
		if err := health.Health(); err != nil {
			logger.Warn().Err(err).Msg("start health failed")
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)

	go es.StartEventsRoutine(cctx, subscriber, &wg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c
	cancel()
	logger.Info().Msg("recv signal, terminating LPTracker ....")
	wg.Wait()
	logger.Info().Msg("LPTracker exit")

	return nil
}
