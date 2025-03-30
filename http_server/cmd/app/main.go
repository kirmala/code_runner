package main

import (
	//"fmt"
	"code_processor/http_server/cmd/app/config"
	"code_processor/http_server/repository/postgres"
	rabbitMQ "code_processor/http_server/repository/rabbit_mq"
	"code_processor/http_server/repository/redis"
	"code_processor/http_server/usecases/service"
	"fmt"
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"code_processor/http_server/api/http"
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

	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	if pgPassword == "" {
		log.Fatal("POSTGRES_PASSWORD is not set")
	}
	pgUser := os.Getenv("POSTGRES_USER")
	if pgUser == "" {
		log.Fatal("POSTGRES_USER is not set")
	}
	pgDB := os.Getenv("POSTGRES_DB")
	if pgDB == "" {
		log.Fatal("POSTGRES_DB is not set")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.Port, pgUser, pgPassword, pgDB)

	taskRepo, err := postgres.NewTaskStorage(connStr)
	if err != nil {
		log.Fatalf("creating task storage: %v", err)
	}
	userRepo, err := postgres.NewUserStorage(connStr)
	if err != nil {
		log.Fatalf("creating user storage: %v", err)
	}

	rdAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	rdPassword := os.Getenv("REDIS_PASSWORD")
	if rdPassword == "" {
		log.Fatal("REDIS_PASSWORD is not set")
	}

	sessionRepo, err := redis.NewSessionStorage(rdAddr, rdPassword)
	if err != nil {
		log.Fatalf("creating session storage: %v", err)
	}
	

	taskSender, err := rabbitMQ.NewRabbitMQSender(rabbitMQAddr, cfg.RabbitMQ.QueueName)
	if err != nil {
		log.Fatalf("failed creating rabbitMQ: %v", err)
	}

	taskService := service.NewTask(taskRepo, sessionRepo, taskSender)
	userService := service.NewUser(userRepo, sessionRepo)

	taskHandlers := http.NewTaskHandler(taskService)
	userHandlers := http.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	taskHandlers.WithTaskHandlers(r)
	userHandlers.WithUserHandlers(r)

	log.Printf("Starting server on %s", addr)
	if err := pkgHttp.CreateAndRunServer(r, addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
