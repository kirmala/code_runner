package main

import (
	"consumer/cmd/app/config"
	"consumer/repository/by_http"
	"consumer/repository/rabbit_mq"
	"consumer/usecases/services/docker_code_processor"
	"fmt"
	"log"
)

// http://localhost:15672/api/healthchecks/node


func main() {
	appFlags := config.ParseFlags()
	var cfg config.AppConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)

	repoAddr := fmt.Sprintf("http://%s:%s/%s", cfg.Repository.Host, cfg.Repository.Port, cfg.Repository.Address)
	taskRepo := byHttp.NewTask(repoAddr)

	imageName := cfg.CodeProcessor.ImageName
	taskCodeProcessor, err := dockerCodeProcessor.NewCodeProcessor(imageName)
	if (err != nil) {
		log.Fatalf("creating docker client: %s", err)
	}

	rabbitMQAddr := fmt.Sprintf("amqp://guest:guest@%s:%s", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)

	TaskReceiver, err := rabbitMQ.NewRabbitMQReceiver(rabbitMQAddr, cfg.RabbitMQ.QueueName, taskCodeProcessor, taskRepo)
	if err != nil {
		log.Fatalf("failed creating rabbitMQ: %v", err)
	}

	TaskReceiver.Receive()
}