package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/pkg/health"
	"starbase.ag/liquidity/pkg/log"
	"starbase.ag/liquidity/pkg/sentry"
	"starbase.ag/liquidity/version"
)

func Run(ctx *cli.Context, module string) *config.Config {
	version.PrintVersion(os.Stdout)

	// logger := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", module)
	// oplog.SetGlobalLogHandler(logger.Handler())

	cfg, err := config.LoadConfigAWS(ctx.String(config.ConfigFlag.Name), config.CONFIG_AWS_ENV_EVENTS)
	if err != nil {
		log.Error("failed to load config", "err", err)
		return nil
	}

	go func() {
		err := health.Health()
		if err != nil {
			log.Error("health check failed", "err", err)
		}
	}()

	if cfg.SentryOpts.DSN != "" {
		_ = sentry.InitSentry(cfg.SentryOpts.DSN, cfg.SentryOpts.Env, cfg.SentryOpts.Debug)
	}

	return &cfg
}
