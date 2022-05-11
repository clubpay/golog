package sentry

import (
	log "github.com/clubpay/golog"
	"time"
)

type config struct {
	dsn          string
	release      string
	env          string
	lvl          log.Level
	tags         map[string]string
	flushTimeout time.Duration
}

type Option func(cfg *config)

func WithRelease(release string) Option {
	return func(cfg *config) {
		cfg.release = release
	}
}

func WithEnv(env string) Option {
	return func(cfg *config) {
		cfg.env = env
	}
}

func WithLevel(level log.Level) Option {
	return func(cfg *config) {
		cfg.lvl = level
	}
}

func WithTags(tags map[string]string) Option {
	return func(cfg *config) {
		cfg.tags = tags
	}
}

func WithFlushTimeout(flushTimeout time.Duration) Option {
	return func(cfg *config) {
		cfg.flushTimeout = flushTimeout
	}
}
