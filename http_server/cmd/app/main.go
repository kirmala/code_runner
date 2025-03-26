package main

import (
	//"fmt"
	"code_processor/http_server/cmd/app/config"
	rabbitMQ "code_processor/http_server/repository/rabbit_mq"
	"code_processor/http_server/repository/ram_storage"
	"code_processor/http_server/usecases/service"
	"fmt"
	"log"

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

	taskRepo := ram_storage.NewTask()
	sessionRepo := ram_storage.NewSession()
	userRepo := ram_storage.NewUser()
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
