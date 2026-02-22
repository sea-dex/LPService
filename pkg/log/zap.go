package log

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	ZapTraceLevel = zapcore.DebugLevel - 1
	zapCallerSkip = 2
)

type zapLogger struct {
	Log        *zap.SugaredLogger
	DesugarLog *zap.Logger
	Options    *Options
}

var _ Logger = (*zapLogger)(nil)

func newZapLogger(opts *Options) Logger {
	cores := make([]zapcore.Core, 0)

	if !opts.Console.Disable {
		w, closeOut, err := zap.Open("stdout")
		if err != nil {
			if closeOut != nil {
				closeOut()
			}
		} else {
			encoder := zapParseEncoder(opts.Console.Encoder, opts.Console.TimeFormat)
			level := ZapParseLevel(opts.Console.Level)
			cores = append(cores, zapcore.NewCore(encoder, w, level))
		}
	}

	if opts.File.Path != "" {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename: opts.File.Path,
			MaxSize:  opts.File.MaxSizeMB,
			MaxAge:   opts.File.MaxAgeDays,
			Compress: true,
		})

		encoder := zapParseEncoder(opts.File.Encoder, opts.File.TimeFormat)
		level := ZapParseLevel(opts.File.Level)
		cores = append(cores, zapcore.NewCore(encoder, w, level))
	}

	if len(cores) == 0 {
		return nil
	}

	l := zap.New(zapcore.NewTee(cores...), zap.WithCaller(true), zap.AddCallerSkip(zapCallerSkip))

	return &zapLogger{
		Log:        l.Sugar(),
		DesugarLog: l,
		Options:    opts,
	}
}

func ZapLogger(l Logger) *zap.Logger {
	if ll, ok := l.(*zapLogger); ok {
		return ll.DesugarLog
	}

	return nil
}

func ZapParseLevel(l string) zapcore.Level {
	switch l {
	case LevelTrace:
		return ZapTraceLevel
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func zapParseEncoder(enc, timeFmt string) zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapParseTimeEncoder(timeFmt)
	cfg.EncodeCaller = zapRelativeCallerEncoder

	switch enc {
	case EncoderJSON:
		cfg.EncodeLevel = zapCapitalLevelEncoder

		return zapcore.NewJSONEncoder(cfg)
	default:
		cfg.EncodeLevel = zapCapitalColorLevelEncoder

		return zapcore.NewConsoleEncoder(cfg)
	}
}

func zapParseTimeEncoder(f string) zapcore.TimeEncoder {
	switch f {
	case "rfc3339nano", "RFC3339Nano":
		return zapcore.RFC3339NanoTimeEncoder
	case "rfc3339", "RFC3339":
		return zapcore.RFC3339TimeEncoder
	case "iso8601", "ISO8601":
		return zapcore.ISO8601TimeEncoder
	case "millis":
		return zapcore.EpochMillisTimeEncoder
	case "nanos":
		return zapcore.EpochNanosTimeEncoder
	case "epoch":
		return zapcore.EpochTimeEncoder
	default:
		return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(f))
		}
	}
}

func zapRelativeCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(zapCallerRelativePath(caller))
}

func zapCapitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l > ZapTraceLevel {
		zapcore.LowercaseLevelEncoder(l, enc)

		return
	}

	enc.AppendString("trace")
}

func zapCapitalColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l > ZapTraceLevel {
		zapcore.CapitalColorLevelEncoder(l, enc)

		return
	}

	enc.AppendString(fmt.Sprintf("\x1b[%dm%s\x1b[0m", color.FgCyan, "TRACE"))
}

func zapCallerRelativePath(ec zapcore.EntryCaller) string {
	if !ec.Defined {
		return "undefined"
	}

	// Find the last separator.
	idx := strings.LastIndexByte(ec.File, '/')
	if idx == -1 {
		return ec.FullPath()
	}

	// Find the penultimate separator.
	idx = strings.LastIndexByte(ec.File[:idx], '/')
	if idx == -1 {
		return ec.FullPath()
	}

	idx = strings.LastIndexByte(ec.Function, '/')
	if idx == -1 {
		return ec.TrimmedPath()
	}

	firstIdx := strings.IndexByte(ec.Function, '/')
	if firstIdx == -1 || firstIdx >= idx {
		return ec.TrimmedPath()
	}

	idx = strings.Index(ec.File, ec.Function[firstIdx+1:idx])
	if idx == -1 {
		return ec.TrimmedPath()
	}

	var sb strings.Builder

	sb.WriteString(ec.File[idx:])
	sb.WriteByte(':')
	sb.WriteString(strconv.Itoa(ec.Line))

	return sb.String()
}

func (l *zapLogger) Tracef(format string, args ...any) {
	l.Log.Log(ZapTraceLevel, fmt.Sprintf(format, args...))
}

func (l *zapLogger) Trace(args ...any) {
	l.Log.Log(ZapTraceLevel, args...)
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.Log.Debugf(format, args...)
}

func (l *zapLogger) Debug(args ...any) {
	l.Log.Debug(args...)
}

func (l *zapLogger) Infof(format string, args ...any) {
	l.Log.Infof(format, args...)
}

func (l *zapLogger) Info(args ...any) {
	l.Log.Info(args...)
}

func (l *zapLogger) Warnf(format string, args ...any) {
	l.Log.Warnf(format, args...)
}

func (l *zapLogger) Warn(args ...any) {
	l.Log.Warn(args...)
}

func (l *zapLogger) Errorf(format string, args ...any) {
	l.Log.Errorf(format, args...)
}

func (l *zapLogger) Error(args ...any) {
	l.Log.Error(args...)
}

func (l *zapLogger) Fatalf(format string, args ...any) {
	l.Log.Fatalf(format, args...)
}

func (l *zapLogger) Fatal(args ...any) {
	l.Log.Fatal(args...)
}

func (l *zapLogger) Panicf(format string, args ...any) {
	l.Log.Fatalf(format, args...)
}

func (l *zapLogger) Panic(args ...any) {
	l.Log.Panic(args...)
}

func (l *zapLogger) Recoverf(err any, format string, args ...any) {
	msg := fmt.Sprintf("[Recovery] "+format, args...)

	l.DesugarLog.Error(msg,
		zap.Time("time", time.Now()),
		zap.String("log-source", "recovery"),
		zap.Any("error", err),
		zap.String("stack", string(debug.Stack())),
	)
}

func (l *zapLogger) Recover(err any, args ...any) {
	msg := fmt.Sprintf("[Recovery]"+strings.Repeat(" %s", len(args)), args...)

	l.DesugarLog.Error(msg,
		zap.Time("time", time.Now()),
		zap.String("log-source", "recovery"),
		zap.Any("error", err),
		zap.String("stack", string(debug.Stack())),
	)
}

func (l *zapLogger) With(fields Fields) Logger {
	f := make([]any, 0)

	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}

	newLogger := l.Log.With(f...)

	return &zapLogger{
		Log:        newLogger,
		DesugarLog: newLogger.Desugar(),
		Options:    l.Options,
	}
}

func (l *zapLogger) WithOption(opts ...zap.Option) Logger {
	newLogger := l.Log.WithOptions(opts...)

	return &zapLogger{
		Log:        newLogger,
		DesugarLog: newLogger.Desugar(),
		Options:    l.Options,
	}
}
