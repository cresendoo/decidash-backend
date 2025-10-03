package xlogger

import (
	"context"
	"log"
	"log/slog"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
	slogmulti "github.com/samber/slog-multi"
)

type Options struct {
	Namespace string
	LogLevel  slog.Level
	SentryDSN string
	AddSource bool
	Release   string
	LogFile    string
	LogFormat  string // "text" 또는 "json"
}

type Option func(*Options)

func WithNamespace(namespace string) Option {
	return func(opts *Options) {
		opts.Namespace = namespace
	}
}

func WithLogLevel(logLevel slog.Level) Option {
	return func(opts *Options) {
		opts.LogLevel = logLevel
	}
}

func WithSentryDSN(sentryDSN string) Option {
	return func(opts *Options) {
		opts.SentryDSN = sentryDSN
	}
}

func WithAddSource(addSource bool) Option {
	return func(opts *Options) {
		opts.AddSource = addSource
	}
}

func WithRelease(release string) Option {
	return func(opts *Options) {
		opts.Release = release
	}
}

func WithLogFile(logFile string) Option {
	return func(opts *Options) {
		opts.LogFile = logFile
	}
}

func WithLogFormat(format string) Option {
	return func(opts *Options) {
		opts.LogFormat = format
	}
}

func Build(ctx context.Context, handlers []slog.Handler, opts ...Option) *slog.Logger {
	o := &Options{AddSource: true, LogFormat: "text"}
	for _, opt := range opts {
		opt(o)
	}

	if o.LogFile != "" {
		fileHandler, err := NewFileHandler(o.LogFile, o.LogFormat, o.LogLevel)
		if err != nil {
			slog.Error("failed to create file handler", "error", err)
		} else {
			handlers = append(handlers, fileHandler)
		}
	}

	if o.SentryDSN != "" {
		if err := sentry.Init(sentry.ClientOptions{Dsn: o.SentryDSN}); err != nil {
			log.Panic(err)
		}
		sentryHandler := sentryslog.Option{Level: slog.LevelError}.NewSentryHandler(ctx)
		handlers = append(handlers, sentryHandler)
	}

	logger := slog.New(slogmulti.Fanout(handlers...)).With("release", o.Release).With("namespace", o.Namespace)
	return logger
}
