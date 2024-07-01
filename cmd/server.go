package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/truongnqse05461/ewallet/internal/api/http/engine"
	"github.com/truongnqse05461/ewallet/internal/api/http/handler"
	"github.com/truongnqse05461/ewallet/internal/api/http/middleware"
	"github.com/truongnqse05461/ewallet/internal/api/http/route"
	"github.com/truongnqse05461/ewallet/internal/config"
	"github.com/truongnqse05461/ewallet/internal/log"
	"github.com/truongnqse05461/ewallet/internal/metrics"
	"github.com/truongnqse05461/ewallet/internal/pg"
	"github.com/truongnqse05461/ewallet/internal/redis"
	"github.com/truongnqse05461/ewallet/internal/service"
	"github.com/truongnqse05461/ewallet/internal/worker"
)

func RunServer() {
	config, err := config.Parse()
	if err != nil {
		panic("config parse failed: " + err.Error())
	}

	logger := log.NewLogger(log.Config{
		Level: config.Log.Level,
	})

	db, err := pg.NewPostgresDB(
		config.Connections.Postgres.Host,
		config.Connections.Postgres.DB,
		config.Connections.Postgres.User,
		config.Connections.Postgres.Password,
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		db.Close()
		logger.Info("database closed")
	}()

	redis, err := redis.NewRedis(
		config.Connections.Redis.Host,
		config.Connections.Redis.Password,
		config.Connections.Redis.WriteTimeout,
		config.Connections.Redis.ReadTimeout,
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		redis.Close()
		logger.Info("redis closed")
	}()

	metric := metrics.New()

	notificationWorker := worker.NewNotificationWorker(logger, 100)

	notificationSvc := service.NewNotificationService(notificationWorker)
	userSvc := service.NewUserService()
	walletSvc := service.NewWalletService()
	transactionSvc := service.NewTransactionService(notificationSvc)

	healthHandler := handler.NewHealthHandler(logger, db)
	swaggerHandler := handler.NewSwaggerHandler(config.Swagger.Host, config.Swagger.BasePath)
	userHandler := handler.NewUserHandler(userSvc)
	walletHandler := handler.NewWalletHandler(walletSvc, transactionSvc)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)

	route := route.NewRoute(
		metric,
		healthHandler,
		swaggerHandler,
		userHandler,
		walletHandler,
		transactionHandler,
	)

	mw := middleware.NewMiddleware(
		logger,
		db,
		metric,
	)

	engine := engine.NewEngine()
	route.Index(engine, mw)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
		Handler: engine,
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.WithErr(err).Error("shutdown server error")
		} else {
			logger.Info("server stopped")
		}
	}()

	defer func() {
		healthHandler.SetShuttingDown()
		logger.Info("set service status to unavailable")
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go notificationWorker.Start(ctx)

	go func() {
		defer close(sigs)

		logger.WithField("port", config.Server.Port).Info("server start")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithErr(err).Error("unexpected stopped")
		}
	}()

	defer func() {
		if sig, ok := <-sigs; ok {
			logger.Infof("received signal: %v", sig)
		}

		logger.Infof("shutting down")
	}()
}
