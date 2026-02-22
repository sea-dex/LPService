package sentry

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"starbase.ag/liquidity/version"
)

func InitSentry(dsn string, env string, debug bool) error {
	opts := sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: dsn,
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug:       debug,
		Release:     version.GitRev,
		Environment: env,
	}

	return sentry.Init(opts)
}

func PanicWithSentry(format string, args ...any) {
	err := fmt.Sprintf(format, args...)
	sentry.CaptureMessage(err)
	sentry.Flush(time.Second * 2)
	panic(err)
}
