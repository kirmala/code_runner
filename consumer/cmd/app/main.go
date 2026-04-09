package main

import (
	"fmt"
	"log"

	"github.com/kirmala/code_runner/consumer/cmd/app/config"
	"github.com/kirmala/code_runner/consumer/internal/api/rabbitmq"
	"github.com/kirmala/code_runner/consumer/internal/repository/postgres"
	"github.com/kirmala/code_runner/consumer/internal/service/basic"
	"github.com/kirmala/code_runner/consumer/internal/service/docker"
)

// http://localhost:15672/api/healthchecks/node

func main() {
	appFlags := config.ParseFlags()
	var cfg config.AppConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.PostgresDB.Host, cfg.PostgresDB.Port, cfg.PostgresDB.User, cfg.PostgresDB.Password, cfg.PostgresDB.DB)
	
	taskRepo, err := postgres.NewTaskStorage(connStr)

	if err != nil {
		log.Fatalf("connecting to postgres: %s", err)
	}

	taskService := basic.NewTask(taskRepo)

	imageName := cfg.ImageName
	clientVersion := cfg.ClientVersion
	runner, err := docker.NewRunner(imageName, clientVersion, cfg.ContainerResource)
	if err != nil {
		log.Fatalf("creating docker client: %s", err)
	}

	rabbitMQAddr := fmt.Sprintf("amqp://guest:guest@%s:%s", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)

	TaskReceiver, err := rabbitmq.NewTaskReceiver(rabbitMQAddr, cfg.QueueName, runner, taskService)
	if err != nil {
		log.Fatalf("failed creating rabbitMQ: %v", err)
	}

	TaskReceiver.Receive()
}
