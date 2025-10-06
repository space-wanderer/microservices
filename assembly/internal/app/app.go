package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/space-wanderer/microservices/assembly/internal/config"
	"github.com/space-wanderer/microservices/platform/pkg/closer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
	apiServer   *http.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	logger.Info(ctx, "🚀 Assembly Service started")

	// Запускаем consumer в горутине
	go func() {
		if err := a.diContainer.ConsumerService(ctx).RunConsumer(ctx); err != nil {
			logger.Error(ctx, "Consumer service error", zap.Error(err))
		}
	}()

	// Запускаем HTTP сервер в горутине
	go func() {
		logger.Info(ctx, "HTTP server listening on :8081")
		if err := a.apiServer.Serve(a.listener); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, "HTTP server error", zap.Error(err))
		}
	}()

	// Ждем завершения контекста
	<-ctx.Done()

	logger.Info(ctx, "🛑 Assembly Service shutting down...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return closer.CloseAll(shutdownCtx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initHTTPServer,
	}
	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(ctx context.Context) error {
	return config.Load()
}

func (a *App) initDi(ctx context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	// Инициализируем логгер с новой конфигурацией
	return logger.InitWithConfig(config.AppConfig().Logger)
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(&logger.NoopLogger{})
	return nil
}

func (a *App) initListener(ctx context.Context) error {
	listener, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		return err
	}

	closer.AddNamed("TCP Listener", func(ctx context.Context) error {
		return listener.Close()
	})

	a.listener = listener
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем базовые middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Добавляем endpoint для метрик
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	// Добавляем health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"status":"healthy"}`)
	})

	a.apiServer = &http.Server{
		Addr:              ":8081",
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	closer.AddNamed("HTTP Server", func(ctx context.Context) error {
		return a.apiServer.Shutdown(ctx)
	})

	return nil
}
