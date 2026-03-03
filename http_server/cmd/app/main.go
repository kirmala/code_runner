package main

import (
	//"fmt"
	"code_processor/http_server/cmd/app/config"
	myenv "code_processor/http_server/cmd/app/env"
	"code_processor/http_server/repository/postgres"
	rabbitMQ "code_processor/http_server/repository/rabbit_mq"
	"code_processor/http_server/repository/redis"
	"code_processor/http_server/service/basic"
	"code_processor/http_server/service/session"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger"

	"code_processor/http_server/api/httpx"
	"code_processor/http_server/api/httpx/middleware"
	_ "code_processor/http_server/docs"
	pkgHttp "code_processor/http_server/pkg/http"
)

// @title code_processor/http_server
// @version 1.0
// @description This is a code runner.

// @host localhost:8080
// @BasePath /
func main() {
	appFlags := config.ParseFlags()
	var cfg config.AppConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)
	addr := fmt.Sprintf("%s:%s", cfg.HTTPConfig.Host, cfg.HTTPConfig.Port)
	rabbitMQAddr := fmt.Sprintf("amqp://guest:guest@%s:%s", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)

	var envCfg myenv.Config
	err := env.Parse(&envCfg)
	if err != nil {
		log.Fatalf("parsing env: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", envCfg.PostgresHost, envCfg.PostgresPort, envCfg.PostgresUser, envCfg.PostgresPassword, envCfg.PostgresDB)
	if err := postgres.RunMigrations(connStr); err != nil {
		log.Fatalf("running migrations: %v", err)
	}

	taskRepo, err := postgres.NewTaskStorage(connStr)
	if err != nil {
		log.Fatalf("creating task storage: %v", err)
	}
	userRepo, err := postgres.NewUserStorage(connStr)
	if err != nil {
		log.Fatalf("creating user storage: %v", err)
	}

	redisCli, err := redis.NewClusterClient(envCfg.RedisAddresses, envCfg.RedisPassword)
	if err != nil {
		log.Fatalf("creating redis client: %v", err)
	}

	sessionRepo := redis.NewSessionStorage(redisCli)

	taskSender, err := rabbitMQ.NewRabbitMQSender(rabbitMQAddr, cfg.QueueName)
	if err != nil {
		log.Fatalf("failed creating rabbitMQ: %v", err)
	}

	taskService := basic.NewTask(taskRepo, sessionRepo, taskSender)
	userService := basic.NewUser(userRepo, sessionRepo)

	taskHandlers := httpx.NewTaskHandler(taskService, session.Authenticator{SessionRepo: sessionRepo})
	userHandlers := httpx.NewUserHandler(userService)

	e := echo.New()
	apiGroup := e.Group("")
	apiGroup.Use(middleware.ServeErrors)
	apiGroup.Use(middleware.Recover)
	taskHandlers.WithTaskHandlers(apiGroup)
	userHandlers.WithUserHandlers(apiGroup)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	log.Printf("Starting server on %s", addr)
	if err := pkgHttp.CreateAndRunServer(e, addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
