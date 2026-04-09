package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/kirmala/code_runner/consumer/cmd/app/config"
	"github.com/kirmala/code_runner/consumer/internal/api/rabbitmq"
	"github.com/kirmala/code_runner/consumer/internal/repository/postgres"
	"github.com/kirmala/code_runner/consumer/internal/service/basic"
	"github.com/kirmala/code_runner/consumer/internal/service/docker"
	slogctx "github.com/veqryn/slog-context"
)

// http://localhost:15672/api/healthchecks/node

func main() {
	h := slogctx.NewHandler(slog.NewJSONHandler(os.Stdout, nil), nil)
	slog.SetDefault(slog.New(h))
	appFlags := config.ParseFlags()
	var cfg config.AppConfig
	err := config.Load(appFlags.ConfigPath, &cfg)
	if err != nil {
		slog.Error("config load failed", slog.Any("error", err))
		return
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.PostgresDB.Host, cfg.PostgresDB.Port, cfg.PostgresDB.User, cfg.PostgresDB.Password, cfg.PostgresDB.DB)

	db, err := postgres.Connect(connStr)
	if err != nil {
		slog.Error("connect postgres failed", slog.Any("error", err))
		return
	}

	taskRepo := postgres.NewTaskStorage(db)

	taskService := basic.NewTask(taskRepo)

	imageName := cfg.ImageName
	clientVersion := cfg.ClientVersion
	runner, err := docker.NewRunner(imageName, clientVersion, cfg.ContainerResource)
	if err != nil {
		slog.Error("create docker client failed", slog.Any("error", err))
		return
	}

	rabbitMQAddr := fmt.Sprintf("amqp://guest:guest@%s:%s", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)

	TaskReceiver, err := rabbitmq.NewTaskReceiver(rabbitMQAddr, cfg.QueueName, runner, taskService)
	if err != nil {
		slog.Error("create rabbitMQ failed", slog.Any("error", err))
	}

	TaskReceiver.Receive()
}
