package log_test

import (
	"testing"
	"time"

	log "github.com/clubpay/golog"
	"github.com/clubpay/golog/core/datadog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLogger(t *testing.T) {
	Convey("Logger", t, func(c C) {
		l := log.New(
			log.WithLevel(log.DebugLevel),
			log.WithJSON(),
			log.WithCore(
				datadog.NewAPI("",
					datadog.WithEnv("test"),
					datadog.WithServiceName("sample"),
					datadog.WithSource("someSource"),
					datadog.WithLevel(log.DebugLevel),
					datadog.WithTags(map[string]string{
						"version": "v1.0.0",
					}),
				),
			),
		)

		l.Debug("Hi", log.String("field1", "f1"))
		time.Sleep(time.Second * 5)
	})
}
