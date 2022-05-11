package datadog

import (
	"context"
	"os"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v2/datadog"
	log "github.com/clubpay/golog"
	"go.uber.org/zap/zapcore"
)

type core struct {
	cfg config
	zapcore.LevelEnabler
	client *datadog.LogsApiService
	enc    log.Encoder
}

func New(apiKey string, opts ...Option) log.Core {
	if apiKey == "" {
		return zapcore.NewNopCore()
	}

	cfg := config{
		flushTimeout: time.Second * 5,
		apiKey:       apiKey,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	_ = os.Setenv("DD_SITE", cfg.site)
	_ = os.Setenv("DD_API_KEY", cfg.apiKey)

	ddConfig := datadog.NewConfiguration()
	ddClient := datadog.NewAPIClient(ddConfig).LogsApi

	return &core{
		cfg:          cfg,
		client:       ddClient,
		LevelEnabler: cfg.lvl,
		enc: log.EncoderBuilder().
			JsonEncoder(),
	}
}

func (c *core) With(fs []log.Field) log.Core {
	return &core{
		cfg:          c.cfg,
		LevelEnabler: c.LevelEnabler,
		client:       c.client,
	}
}

func (c *core) Check(ent log.Entry, ce *log.CheckedEntry) *log.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}

	return ce
}

func (c *core) Write(ent log.Entry, fs []log.Field) error {
	m := make(map[string]interface{}, len(fs))
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fs {
		f.AddTo(enc)
	}
	for k, v := range enc.Fields {
		m[k] = v
	}

	buf, err := c.enc.EncodeEntry(ent, fs)
	if err != nil {
		return err
	}
	body := []datadog.HTTPLogItem{
		{
			Ddsource: c.cfg.source,
			Ddtags:   c.cfg.tags,
			Hostname: c.cfg.hostname,
			Message:  datadog.PtrString(buf.String()),
			Service:  c.cfg.service,
		},
	}

	ctx, cf := context.WithTimeout(context.Background(), c.cfg.flushTimeout)
	defer cf()
	_, _, err = c.client.SubmitLog(
		datadog.NewDefaultContext(ctx),
		body,
		*datadog.NewSubmitLogOptionalParameters().
			WithContentEncoding(datadog.CONTENTENCODING_DEFLATE),
	)

	buf.Free()

	return err
}

func (c *core) Sync() error {
	return nil
}
