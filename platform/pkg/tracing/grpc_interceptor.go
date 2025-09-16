package tracing

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor создает unary client interceptor для трейсинга gRPC клиентов
func UnaryClientInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)
}

// UnaryServerInterceptor создает unary server interceptor для трейсинга gRPC серверов
func UnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)
}

// StreamClientInterceptor создает stream client interceptor для трейсинга gRPC клиентов
func StreamClientInterceptor(serviceName string) grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)
}

// StreamServerInterceptor создает stream server interceptor для трейсинга gRPC серверов
func StreamServerInterceptor(serviceName string) grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)
}

// StartSpan создает новый span с указанным именем
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	tracer := otel.Tracer("microservices")
	return tracer.Start(ctx, name, opts...)
}
