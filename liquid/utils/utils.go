package utils

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"starbase.ag/liquidity/version"
)

// InitSentry initialize sentry.
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

// PanicWithSentry sentry then panic.
func PanicWithSentry(format string, args ...any) {
	err := fmt.Sprintf(format, args...)
	sentry.CaptureMessage(err)
	sentry.Flush(time.Second * 2)

	log.Fatal().Stack().Msgf(format, args...)
}

func HashToAddress(hash common.Hash) string {
	addr := common.Address{}
	addr.SetBytes(hash.Bytes())

	return addr.Hex()
}

func GetBoolVariableParam(opts ...any) bool {
	wait := false

	if len(opts) > 0 {
		if v, ok := opts[0].(bool); ok {
			wait = v
		}
	}

	return wait
}

func ToBigIntMust(s string) *big.Int {
	v, ok := big.NewInt(0).SetString(s, 0)
	if !ok {
		panic("convert string to bigInt failed")
	}

	return v
}
