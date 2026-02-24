package rabbitmq

import (
	"code_processor/consumer/service"
	"code_processor/http_server/models"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type TaskReceiver struct {
	connection  *amqp.Connection
	channel     *amqp.Channel
	queueName   string
	runner      service.Runner
	taskService service.Task
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func NewTaskReceiver(ampqURL, queueName string, runner service.Runner, taskService service.Task) (*TaskReceiver, error) {
	conn, err := amqp.Dial(ampqURL)
	if err != nil {
		return nil, fmt.Errorf("connecting to rabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("creating channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue declaration: %w", err)
	}

	return &TaskReceiver{
		connection:  conn,
		channel:     ch,
		queueName:   queueName,
		runner:      runner,
		taskService: taskService,
	}, nil
}

func (r *TaskReceiver) Receive() {
	msgs, err := r.channel.Consume(
		r.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var task models.Task
			errorTask := task
			err := json.Unmarshal([]byte(d.Body), &task)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				errorTask.Result = err.Error()
				errorTask.Status = models.StatusFailed

				err = r.taskService.Put(context.Background(), errorTask)
				if err != nil {
					log.Printf("adding processed task to database: %s", err)
				}
				continue
			}
			processedTask, err := r.runner.Run(context.Background(), task)
			if err != nil {
				log.Printf("processing task: %s", err)
				errorTask.Result = err.Error()
				errorTask.Status = models.StatusFailed
				err = r.taskService.Put(context.Background(), errorTask)
				if err != nil {
					log.Printf("adding processed task to database: %s", err)
				}
				continue
			}
			err = r.taskService.Put(context.Background(), processedTask)
			if err != nil {
				log.Printf("adding processed task to database: %s", err)
				continue
			}
			log.Printf("Success")
		}
	}()
	<-forever
}

func (r *TaskReceiver) Close() error {
	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("closing ampq channel: %w", err)
	}
	if err := r.connection.Close(); err != nil {
		return fmt.Errorf("closing ampq connection: %w", err)
	}
	return nil
}
