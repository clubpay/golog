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

func (cfg *config) tagsToStr() {
	if cfg.env != nil {
		cfg.tags["env"] = *cfg.env
	}
	if cfg.source != nil {
		cfg.tags["source"] = *cfg.source
	}
	if cfg.service != nil {
		cfg.tags["service"] = *cfg.service
	}
	if cfg.hostname != nil {
		cfg.tags["hostname"] = *cfg.hostname
	}
	if len(cfg.tags) > 0 {
		sb := strings.Builder{}
		idx := 0
		for k, v := range cfg.tags {
			if idx != 0 {
				sb.WriteString(",")
			}
			sb.WriteString(k)
			sb.WriteString(":")
			sb.WriteString(v)
			idx++
		}
		cfg.tagsStr = datadog.PtrString(sb.String())
	}
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
