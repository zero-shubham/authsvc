// Package config provides initialization of OpenTelemetry setup.
package config // import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/example/config"

import (
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Init configures an OpenTelemetry exporter and trace provider.
func Init() (*sdktrace.TracerProvider, error) {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	processor := sdktrace.NewBatchSpanProcessor(exporter)
	tp.RegisterSpanProcessor(processor)

	return tp, nil

}

// https://pkg.go.dev/go.opentelemetry.io/otel/metric#example-Meter-Asynchronous_multiple
