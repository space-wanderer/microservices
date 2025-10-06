package tracing

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

// NewClientStatsHandler создает stats handler для трейсинга gRPC клиентов
func NewClientStatsHandler(serviceName string) stats.Handler {
	return otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)
}

// NewServerStatsHandler создает stats handler для трейсинга gRPC серверов
func NewServerStatsHandler(serviceName string) stats.Handler {
	return otelgrpc.NewServerHandler(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)
}

// UnaryClientInterceptor создает unary client interceptor для трейсинга gRPC клиентов
// Deprecated: используйте NewClientStatsHandler с grpc.WithStatsHandler
func UnaryClientInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	// Возвращаем простой interceptor, который не делает ничего
	// В новой версии рекомендуется использовать stats.Handler
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// UnaryServerInterceptor создает unary server interceptor для трейсинга gRPC серверов
// Deprecated: используйте NewServerStatsHandler с grpc.StatsHandler
func UnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	// Возвращаем простой interceptor, который не делает ничего
	// В новой версии рекомендуется использовать stats.Handler
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		return handler(ctx, req)
	}
}
