package log

import (
	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

type SentryOptionFunc func(o *zapsentry.Configuration)

func AttachSentryCore(opts ...SentryOptionFunc) {
	ll, ok := log.(*zapLogger)
	if !ok || log == nil {
		return
	}

	clt := sentry.CurrentHub().Client()
	if clt == nil {
		return
	}

	conf := &zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, // when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.DebugLevel, // at what level should we sent breadcrumbs to sentry
	}

	for _, opt := range opts {
		opt(conf)
	}

	core, err := zapsentry.NewCore(*conf, zapsentry.NewSentryClientFromClient(clt))
	if err != nil {
		return
	}

	al := zapsentry.AttachCoreToLogger(core, ll.DesugarLog)

	ll.Log = al.Sugar()
	ll.DesugarLog = al

	log = ll
}

func WithSentrySkipModulePrefixes(ms ...string) SentryOptionFunc {
	return func(o *zapsentry.Configuration) {
		skips := make([]zapsentry.FrameMatcher, 0, len(ms))
		for _, m := range ms {
			skips = append(skips, zapsentry.SkipModulePrefixFrameMatcher(m))
		}

		o.FrameMatcher = zapsentry.CombineFrameMatchers(skips...)
	}
}

func WithSentryDisableStacktrace(x bool) SentryOptionFunc {
	return func(o *zapsentry.Configuration) {
		o.DisableStacktrace = x
	}
}
