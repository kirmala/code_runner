package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/kirmala/code_runner/http_server/cmd/app/config"
	"github.com/kirmala/code_runner/http_server/internal/metrics"
	"github.com/kirmala/code_runner/http_server/internal/repository/postgres"
	rabbitMQ "github.com/kirmala/code_runner/http_server/internal/repository/rabbit_mq"
	"github.com/kirmala/code_runner/http_server/internal/repository/redis"
	"github.com/kirmala/code_runner/http_server/internal/service/basic"
	"github.com/kirmala/code_runner/http_server/internal/service/session"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"

	_ "github.com/kirmala/code_runner/http_server/docs"
	"github.com/kirmala/code_runner/http_server/internal/api/httpx"
	"github.com/kirmala/code_runner/http_server/internal/api/httpx/middleware"
	pkgHttp "github.com/kirmala/code_runner/http_server/pkg/http"
	slogctx "github.com/veqryn/slog-context"
)

// @title github.com/kirmala/code_runner/http_server
// @version 1.0
// @description This is a code runner.

// @host localhost:8080
// @BasePath /
func main() {
	h := slogctx.NewHandler(slog.NewJSONHandler(os.Stdout, nil), nil)
	slog.SetDefault(slog.New(h).With(slog.String("service", "api")))

	appFlags := config.ParseFlags()
	var cfg config.AppConfig
	config.Load(appFlags.ConfigPath, &cfg)
	
	addr := fmt.Sprintf("%s:%s", cfg.HTTPConfig.Host, cfg.HTTPConfig.Port)
	rabbitMQAddr := fmt.Sprintf("amqp://guest:guest@%s:%s", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.PostgresDB.Host, cfg.PostgresDB.Port, cfg.PostgresDB.User, cfg.PostgresDB.Password, cfg.PostgresDB.DB)

	db, err := postgres.Connect(connStr)

	if err := postgres.RunMigrations(db); err != nil {
		slog.Error("failed to run migrations", slog.Any("error", err))
		return
	}
	taskRepo := postgres.NewTaskStorage(db)
	userRepo := postgres.NewUserStorage(db)

	redisCli, err := redis.NewClusterClient(cfg.RedisDB.Addresses, cfg.RedisDB.Password)
	if err != nil {
		slog.Error("failed to create redis client", slog.Any("error", err))
		return
	}

	sessionRepo := redis.NewSessionStorage(redisCli)

	taskSender, err := rabbitMQ.NewRabbitMQSender(rabbitMQAddr, cfg.QueueName)
	if err != nil {
		slog.Error("failed to create rabbitmq", slog.Any("error", err))
		return
	}

	taskService := basic.NewTask(taskRepo, sessionRepo, taskSender)
	userService := basic.NewUser(userRepo, sessionRepo)

	taskHandlers := httpx.NewTaskHandler(taskService, session.Authenticator{SessionRepo: sessionRepo})
	userHandlers := httpx.NewUserHandler(userService)

	e := echo.New()
	apiGroup := e.Group("")

	apiGroup.Use(middleware.CorrelationID)
	apiGroup.Use(middleware.Metrics)
	apiGroup.Use(middleware.ServeErrors)
	apiGroup.Use(middleware.Logger)
	apiGroup.Use(middleware.Recover)
	
	taskHandlers.WithTaskHandlers(apiGroup)
	userHandlers.WithUserHandlers(apiGroup)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	reg := prometheus.NewRegistry()
	metrics.Register(reg)

	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		http.ListenAndServe(":2112", nil)
	}()

	slog.Info("starting server", slog.String("address", addr))
	if err := pkgHttp.CreateAndRunServer(e, addr); err != nil {
		slog.Error("failed to start server", slog.Any("error", err))
	}
}
