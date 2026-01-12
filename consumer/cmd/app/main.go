package main

import (
	"code_processor/consumer/api/rabbit_mq"
	"code_processor/consumer/cmd/app/config"
	//"code_processor/consumer/repository/by_http"
	"code_processor/consumer/repository/postgres"
	"code_processor/consumer/usecases/services/docker_code_processor"
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
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.Port, pgUser, pgPassword, pgDB)
	taskRepo, err := postgres.NewTaskStorage(connStr)

	if err != nil {
		log.Fatalf("connecting to postgres: %s", err)
	}

	imageName := cfg.CodeProcessor.ImageName
	taskCodeProcessor, err := dockerCodeProcessor.NewCodeProcessor(imageName)
	if err != nil {
		log.Fatalf("creating docker client: %s", err)
	}

	rabbitMQAddr := fmt.Sprintf("amqp://guest:guest@%s:%s", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)

	TaskReceiver, err := rabbitMQ.NewRabbitMQReceiver(rabbitMQAddr, cfg.RabbitMQ.QueueName, taskCodeProcessor, taskRepo)
	if err != nil {
		log.Fatalf("failed creating rabbitMQ: %v", err)
	}

	TaskReceiver.Receive()
}
