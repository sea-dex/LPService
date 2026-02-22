package log

const (
	LevelTrace = "trace"
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)

const (
	EncoderConsole = "console"
	EncoderJSON    = "json"
)

var log Logger

type Fields map[string]any

type Logger interface {
	Tracef(format string, args ...any)
	Trace(args ...any)

	Debugf(format string, args ...any)
	Debug(args ...any)

	Infof(format string, args ...any)
	Info(args ...any)

	Warnf(format string, args ...any)
	Warn(args ...any)

	Errorf(format string, args ...any)
	Error(args ...any)

	Fatalf(format string, args ...any)
	Fatal(args ...any)

	Panicf(format string, args ...any)
	Panic(args ...any)

	Recoverf(err any, format string, args ...any)
	Recover(err any, args ...any)

	With(fields Fields) Logger
}

func InitLogger(opts ...Option) {
	log = NewLogger(opts...)
}

func NewLogger(opts ...Option) Logger {
	return newZapLogger(newOptions(opts...))
}

func Tracef(format string, args ...any) {
	log.Tracef(format, args...)
}

func Trace(args ...any) {
	log.Trace(args...)
}

func Debugf(format string, args ...any) {
	log.Debugf(format, args...)
}

func Debug(args ...any) {
	log.Debug(args...)
}

func Infof(format string, args ...any) {
	log.Infof(format, args...)
}

func Info(args ...any) {
	log.Info(args...)
}

func Warnf(format string, args ...any) {
	log.Warnf(format, args...)
}

func Warn(args ...any) {
	log.Warn(args...)
}

func Errorf(format string, args ...any) {
	args = appendStackTraceMaybeArgs(args)
	log.Errorf(format, args...)
}

func Error(args ...any) {
	args = appendStackTraceMaybeArgs(args)
	log.Error(args...)
}

func Fatalf(format string, args ...any) {
	args = appendStackTraceMaybeArgs(args)
	log.Fatalf(format, args...)
}

func Fatal(args ...any) {
	args = appendStackTraceMaybeArgs(args)
	log.Fatal(args...)
}

func Panicf(format string, args ...any) {
	log.Panicf(format, args...)
}

func Panic(args ...any) {
	log.Panic(args...)
}

func Recoverf(err any, format string, args ...any) {
	log.Recoverf(err, format, args...)
}

func Recover(err any, args ...any) {
	log.Recover(err, args...)
}

func With(keyValues Fields) Logger {
	return log.With(keyValues)
}
