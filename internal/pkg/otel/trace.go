package otel

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type TraceProviderOption struct {
	AppName string
}

func NewTraceExporter(logger log.Logger, c *conf.Observability) *otlptrace.Exporter {
	conf := c.Trace
	opts := []otlptracegrpc.Option{}

	if conf.Endpoint != "" {
		opts = append(opts, otlptracegrpc.WithEndpoint(conf.Endpoint), otlptracegrpc.WithInsecure())
	}
	exp, err := otlptracegrpc.New(context.Background(), opts...)
	if err != nil {
		log.NewHelper(logger).Fatalf("failed to initialize trace exporter: %v", err)
		return nil
	}
	return exp
}

func NewSampler(fraction float64) tracesdk.Sampler {
	if fraction < 0 {
		return tracesdk.ParentBased(tracesdk.AlwaysSample())
	}
	return tracesdk.ParentBased(tracesdk.TraceIDRatioBased(fraction))
}

func InitTraceProvider(logger log.Logger, sampler tracesdk.Sampler, exp *otlptrace.Exporter, res *resourcesdk.Resource) ShutdownFunc {
	// todo with sampler
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(sampler),
	)
	otel.SetTracerProvider(tp)
	return func(ctx context.Context) {
		log.NewHelper(logger).Infof("start shutdown trace provider...")
		if err := tp.Shutdown(ctx); err != nil {
			log.NewHelper(logger).Errorf("failed to shutdown trace provider: %v", err)
		}
	}
}
