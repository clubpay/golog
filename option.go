package log

import "go.uber.org/zap/zapcore"

type Option func(cfg *config)

type config struct {
	level           Level
	TimeEncoder     TimeEncoder
	LevelEncoder    LevelEncoder
	DurationEncoder DurationEncoder
	CallerEncoder   CallerEncoder
	skipCaller      int
	encoder         string

	cores []Core
}

var defaultConfig = config{
	level:           InfoLevel,
	skipCaller:      1,
	encoder:         "console",
	TimeEncoder:     timeEncoder,
	LevelEncoder:    zapcore.CapitalLevelEncoder,
	DurationEncoder: zapcore.StringDurationEncoder,
	CallerEncoder:   zapcore.ShortCallerEncoder,
}

func WithLevel(lvl Level) Option {
	return func(cfg *config) {
		cfg.level = lvl
	}
}

func WithSkipCaller(skip int) Option {
	return func(cfg *config) {
		cfg.skipCaller = skip
	}
}

func WithJSON() Option {
	return func(cfg *config) {
		cfg.encoder = "json"
	}
}

func WithConsole() Option {
	return func(cfg *config) {
		cfg.encoder = "console"
	}
}

func WithCore(core Core) Option {
	return func(cfg *config) {
		cfg.cores = append(cfg.cores, core)
	}
}
