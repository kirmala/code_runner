package rabbitMQ

import (
	"code_processor/consumer/repository"
	"code_processor/consumer/usecases"
	"code_processor/http_server/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQReceiver struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	queueName     string
	codeProcessor usecases.CodeProcessor
	repo          repository.Task
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func NewRabbitMQReceiver(ampqURL, queueName string, codeProcessor usecases.CodeProcessor, repo repository.Task) (*RabbitMQReceiver, error) {
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

	return &RabbitMQReceiver{
		connection:    conn,
		channel:       ch,
		queueName:     queueName,
		codeProcessor: codeProcessor,
		repo:          repo,
	}, nil
}

func (r *RabbitMQReceiver) Receive() {
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
				errorTask.Result = "ready"
				errorTask.Status = err.Error()
				err = r.repo.Put(errorTask)
				if err != nil {
					log.Printf("adding processed task to database: %s", err)
				}
				continue
			}
			processedTask, err := r.codeProcessor.Process(task)
			if err != nil {
				log.Printf("processing task: %s", err)
				errorTask.Result = "ready"
				errorTask.Status = err.Error()
				err = r.repo.Put(errorTask)
				if err != nil {
					log.Printf("adding processed task to database: %s", err)
				}
				continue
			}
			err = r.repo.Put(*processedTask)
			if err != nil {
				log.Printf("adding processed task to database: %s", err)
				continue
			}
			log.Printf("Success")
		}
	}()
	<-forever
}

func (r *RabbitMQReceiver) Close() {
	r.channel.Close()
	r.connection.Close()
}
