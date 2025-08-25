package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/space-wanderer/microservices/order/internal/config"
	"github.com/space-wanderer/microservices/platform/pkg/closer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
	order_v1 "github.com/space-wanderer/microservices/shared/pkg/api/order/v1"
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
	// Запускаем Kafka consumer в горутине
	go func() {
		consumerService := a.diContainer.ShipAssembledConsumerService(ctx)
		if consumerService != nil {
			if err := consumerService.RunConsumer(ctx); err != nil {
				logger.Error(ctx, "Failed to run Kafka consumer", zap.Error(err))
			}
		}
	}()

	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initMigrations,
		a.initHTTPServer,
	}
	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDi(ctx context.Context) error {
	a.diContainer = NewDiContainer()

	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())

	return nil
}

func (a *App) initListener(ctx context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().OrderHTTP.Address())
	if err != nil {
		return err
	}

	closer.AddNamed("TCP Listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	// Создаем OpenAPI сервер
	api := a.diContainer.OrderV1API(ctx)
	s, err := order_v1.NewServer(api)
	if err != nil {
		return fmt.Errorf("failed to create OpenAPI server: %w", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", s)

	a.apiServer = &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	closer.AddNamed("HTTP Server", func(ctx context.Context) error {
		return a.apiServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) initMigrations(ctx context.Context) error {
	migrator := a.diContainer.PGMigrator(ctx)
	if migrator == nil {
		return fmt.Errorf("failed to create migrator")
	}

	err := migrator.Up()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("HTTP server listening on %s", config.AppConfig().OrderHTTP.Address()))

	err := a.apiServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) HTTPServer() *http.Server {
	return a.apiServer
}
