package datadog

import (
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v2/datadog"

	log "github.com/clubpay/golog"
)

const (
	US = "datadoghq.us"
	EU = "datadoghq.eu"
)

type config struct {
	apiKey        string
	agentHostPort string
	site          string
	lvl           log.Level
	flushTimeout  time.Duration

	tags     map[string]string
	tagsStr  *string
	env      *string
	source   *string
	hostname *string
	service  *string
}

type Option func(cfg *config)

func WithLevel(level log.Level) Option {
	return func(cfg *config) {
		cfg.lvl = level
	}
}

func WithFlushTimeout(flushTimeout time.Duration) Option {
	return func(cfg *config) {
		cfg.flushTimeout = flushTimeout
	}
}

func WithTags(tags map[string]string) Option {
	return func(cfg *config) {
		cfg.tags = tags
	}
}

func WithServiceName(serviceName string) Option {
	return func(cfg *config) {
		cfg.service = datadog.PtrString(serviceName)
	}
}

func WithSource(source string) Option {
	return func(cfg *config) {
		cfg.source = datadog.PtrString(source)
	}
}

func WithEnv(env string) Option {
	return func(cfg *config) {
		cfg.env = datadog.PtrString(env)
	}
}

func WithHostname(hostname string) Option {
	return func(cfg *config) {
		cfg.hostname = datadog.PtrString(hostname)
	}
}

func WithSite(site string) Option {
	return func(cfg *config) {
		cfg.site = site
	}
}
