package logger

import (
	"io"
	"os"
	"path"

	zlogsentry "github.com/archdx/zerolog-sentry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"starbase.ag/liquidity/config"
	"starbase.ag/liquidity/version"
)

var Logger zerolog.Logger

// Init initialize logger with sentry.
func Init(dsn string, env string, debug bool, args ...any) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	tmFmt := "2006-01-02T15:04:05.000Z07:00"
	writers := []io.Writer{}

	if env == "local" {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: tmFmt})
	} else {
		writers = append(writers, os.Stdout)
	}

	var conf *config.Log
	if len(args) > 0 {
		conf = args[0].(*config.Log)
	}

	if conf != nil && conf.FileLoggingEnabled {
		w := newRollingFile(conf)
		if w != nil {
			writers = append(writers, w)
		}
	}

	if dsn != "" {
		w, err := zlogsentry.New(dsn,
			zlogsentry.WithEnvironment(env),
			zlogsentry.WithRelease(version.Version),
			zlogsentry.WithLevels(zerolog.ErrorLevel, zerolog.FatalLevel),
		)
		if err != nil {
			panic(err.Error())
		}

		// defer w.Close()
		writers = append(writers, w)
	}
	// else {
	// 	if env == "local" {
	// 		// append log to file
	// 		runLogFile, _ := os.OpenFile(
	// 			"lpservice.log",
	// 			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
	// 			0o664,
	// 		)
	// 		multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: tmFmt}, runLogFile)
	// 		Logger = zerolog.New(multi).With().Timestamp().Logger()
	// 	} else {
	// 		if debug {
	// 			Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: tmFmt}).With().Timestamp().Logger()
	// 		} else {
	// 			Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	// 		}
	// 	}
	// }

	multi := zerolog.MultiLevelWriter(writers...)
	Logger = zerolog.New(multi).With().Timestamp().Logger()
}

func newRollingFile(lc *config.Log) io.Writer {
	if err := os.MkdirAll(lc.Directory, 0o744); err != nil {
		log.Error().Err(err).Str("path", lc.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(lc.Directory, lc.Filename),
		MaxBackups: lc.MaxBackups, // files
		MaxSize:    lc.MaxSize,    // megabytes
		MaxAge:     lc.MaxAge,     // days
	}
}

// Info info level.
func Info() *zerolog.Event {
	return Logger.Info()
}

// Warn warn level.
func Warn() *zerolog.Event {
	return Logger.Warn()
}

// Error error level.
func Error() *zerolog.Event {
	return Logger.Error()
}

// Fatal fatal level.
func Fatal() *zerolog.Event {
	return Logger.Fatal()
}
