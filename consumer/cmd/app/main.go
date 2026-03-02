package main

import (
	"code_processor/consumer/api/rabbitmq"
	"code_processor/consumer/cmd/app/config"
	"code_processor/consumer/service/basic"
	"code_processor/consumer/service/docker"

	//"code_processor/consumer/repository/by_http"
	"code_processor/consumer/repository/postgres"
	"os"

	"fmt"
	"log"
)

// http://localhost:15672/api/healthchecks/node

func main() {
	appFlags := config.ParseFlags()
	var cfg config.AppConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)

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
	pgHost := os.Getenv("POSTGRES_HOST")
	if pgHost == "" {
		log.Fatal("POSTGRES_HOST is not set")
	}
	pgPort := os.Getenv("POSTGRES_PORT")
	if pgPort == "" {
		log.Fatal("POSTGRES_PORT is not set")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", pgHost, pgPort, pgUser, pgPassword, pgDB)
	
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
