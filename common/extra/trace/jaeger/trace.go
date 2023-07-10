package trace

import (
	"errors"
	"io"
	"time"

	pconfig "github.com/star-table/common/core/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// NewTracer 创建一个jaeger Tracer
func NewTracer(conf *pconfig.JaegerConfig) (opentracing.Tracer, io.Closer, error) {
	if conf == nil{
		return nil, nil, errors.New("conf is nil")
	}
	if conf.SamplerType == ""{
		conf.SamplerType = jaeger.SamplerTypeProbabilistic
	}

	cfg := config.Configuration{
		ServiceName: conf.TraceService,
		Sampler: &config.SamplerConfig{
			Type:  conf.SamplerType,
			Param: conf.SamplerParam,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	sender, err := jaeger.NewUDPTransport(conf.UdpAddress, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		config.Reporter(reporter),
	)

	return tracer, closer, err
}

