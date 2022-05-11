package datadog

import (
	"strings"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v2/datadog"

	log "github.com/clubpay/golog"
)

const (
	US = "datadoghq.us"
	EU = "datadoghq.eu"
)

type config struct {
	apiKey       string
	site         string
	lvl          log.Level
	flushTimeout time.Duration

	tags     *string
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
		sb := strings.Builder{}
		idx := 0
		for k, v := range tags {
			if idx != 0 {
				sb.WriteString(",")
			}
			sb.WriteString(k)
			sb.WriteString(":")
			sb.WriteString(v)
			idx++
		}

		cfg.tags = datadog.PtrString(sb.String())
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
