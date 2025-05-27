// Package config provides initialization of OpenTelemetry setup.
package config // import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/example/config"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otlpmetrichttp "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	otlptracehttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmeter "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Init configures an OpenTelemetry exporter and trace provider.
func Init(ctx context.Context, otelEndpoint string) (*sdktrace.TracerProvider, *sdkmeter.MeterProvider, error) {
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otelEndpoint),
		otlptracehttp.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	metricExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(otelEndpoint),
		otlpmetrichttp.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	resource := resource.NewWithAttributes(
		"authsvc",
		attribute.String("service.name", "authsvc"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	mp := sdkmeter.NewMeterProvider(
		sdkmeter.WithReader(sdkmeter.NewPeriodicReader(metricExporter, sdkmeter.WithInterval(time.Second*3))),
		sdkmeter.WithResource(resource),
	)
	otel.SetMeterProvider(mp)

	processor := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp.RegisterSpanProcessor(processor)

	return tp, mp, nil

}

// https://pkg.go.dev/go.opentelemetry.io/otel/metric#example-Meter-Asynchronous_multiple
