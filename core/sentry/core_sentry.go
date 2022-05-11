package sentry

import (
	log "github.com/clubpay/golog"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

/*
   Creation Time: 2019 - Apr - 24
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type sentryCore struct {
	zapcore.LevelEnabler
	hub          *sentry.Hub
	tags         map[string]string
	flushTimeout time.Duration
}

func New(dsn string, opts ...Option) log.Core {
	cfg := config{
		flushTimeout: time.Second * 5,
		dsn:          dsn,
		lvl:          log.WarnLevel,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         cfg.dsn,
		Release:     cfg.release,
		Environment: cfg.env,
	})
	if err != nil {
		return zapcore.NewNopCore()
	}

	sentryScope := sentry.NewScope()
	sentryHub := sentry.NewHub(client, sentryScope)

	return &sentryCore{
		hub:          sentryHub,
		tags:         cfg.tags,
		LevelEnabler: cfg.lvl,
		flushTimeout: cfg.flushTimeout,
	}
}

func (c *sentryCore) With(fs []log.Field) log.Core {
	return &sentryCore{
		hub:          c.hub,
		tags:         c.tags,
		LevelEnabler: c.LevelEnabler,
	}
}

func (c *sentryCore) Check(ent log.Entry, ce *log.CheckedEntry) *log.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}

	return ce
}

func (c *sentryCore) Write(ent log.Entry, fs []log.Field) error {
	m := make(map[string]interface{}, len(fs))
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fs {
		f.AddTo(enc)
	}
	for k, v := range enc.Fields {
		m[k] = v
	}

	event := sentry.NewEvent()
	event.Message = ent.Message
	event.Timestamp = ent.Time
	event.Level = sentryLevel(ent.Level)
	event.Extra = m
	event.Tags = c.tags
	c.hub.CaptureEvent(event)

	// We may be crashing the program, so should flush any buffered events.
	if ent.Level > log.ErrorLevel {
		c.hub.Flush(time.Second)
	}

	return nil
}

func (c *sentryCore) Sync() error {
	c.hub.Flush(c.flushTimeout)

	return nil
}

func sentryLevel(lvl log.Level) sentry.Level {
	switch lvl {
	case log.DebugLevel:
		return sentry.LevelDebug
	case log.InfoLevel:
		return sentry.LevelInfo
	case log.WarnLevel:
		return sentry.LevelWarning
	case log.ErrorLevel:
		return sentry.LevelError
	case log.PanicLevel:
		return sentry.LevelFatal
	case log.FatalLevel:
		return sentry.LevelFatal
	default:
		// Unrecognized levels are fatal.
		return sentry.LevelFatal
	}
}
