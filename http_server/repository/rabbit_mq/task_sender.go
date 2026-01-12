package rabbitMQ

import (
	"code_processor/http_server/models"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQSender struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

func NewRabbitMQSender(ampqURL, queueName string) (*RabbitMQSender, error) {
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

	return &RabbitMQSender{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
	}, nil
}

func (r *RabbitMQSender) Send(task models.Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("task mashal: %w", err)
	}

	err = r.channel.Publish(
		"",
		r.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("task publishing: %w", err)
	}

	return nil
}

func (r *RabbitMQSender) Close() error {
	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("closing ampq channel: %w", err)
	}
	if err := r.connection.Close(); err != nil {
		return fmt.Errorf("closing ampq connection: %w", err)
	}
	return nil
}
