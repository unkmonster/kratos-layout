package otel

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/attribute"
	resourcesdk "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type ShutdownFunc func(context.Context)

func NewResource(logger log.Logger, conf *conf.App, appName string) *resourcesdk.Resource {
	res, err := resourcesdk.New(
		context.Background(),
		resourcesdk.WithOS(),
		resourcesdk.WithContainer(),
		resourcesdk.WithContainerID(),
		resourcesdk.WithHost(),
		resourcesdk.WithTelemetrySDK(),
		resourcesdk.WithAttributes(
			attribute.String("env", conf.Env),
			semconv.ServiceNameKey.String(appName),
		),
	)
	if err != nil {
		log.NewHelper(logger).Fatalf("failed to new resource: %v", err)
		return nil
	}

	res, err = resourcesdk.Merge(resourcesdk.Default(), res)
	if err != nil {
		log.NewHelper(logger).Fatalf("failed to merge resource: %v", err)
		return nil
	}
	return res
}
