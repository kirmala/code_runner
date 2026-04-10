package rabbitMQ

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/proto"

	"github.com/kirmala/code_runner/contracts/gen/pb"
	"github.com/kirmala/code_runner/http_server/internal/domain"
	"github.com/kirmala/code_runner/http_server/pkg/correlationid"
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

func (r *RabbitMQSender) Send(ctx context.Context, task domain.Task) error {
	msg := pb.TaskExecutionMessage{
		TaskId: task.Id.String(),
		Code: task.Code,
		Translator: pb.TaskTranslator(task.Translator),
	}
	body, err := proto.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("task mashal: %w", err)
	}

	correlationId, ok := correlationid.FromContext(ctx)
	if !ok {
		slog.WarnContext(ctx, "No correlationId in context")
	}

	err = r.channel.Publish(
		"",
		r.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/x-protobuf",
			Body:        body,
			CorrelationId: correlationId,
		},
	)
	if err != nil {
		return fmt.Errorf("task publishing: %w", err)
	}

	slog.InfoContext(ctx, "task sent successfully")

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
