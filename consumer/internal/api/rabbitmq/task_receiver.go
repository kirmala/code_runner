package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/consumer/internal/domain"
	"github.com/kirmala/code_runner/consumer/internal/service"
	"github.com/kirmala/code_runner/contracts/gen/pb"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type TaskReceiver struct {
	connection  *amqp.Connection
	channel     *amqp.Channel
	queueName   string
	runner      service.Runner
	taskService service.Task

	// middlewares to be applied before running handler
	// where 0 is the closest to handler middleware
	middlewares []Middleware
}

func NewTaskReceiver(ampqURL, queueName string, runner service.Runner, taskService service.Task, middlewares []Middleware) (*TaskReceiver, error) {
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
		middlewares: middlewares,
	}, nil
}

// Handler process rabbitmq message
type Handler func(context.Context, amqp.Delivery) error

// A Middleware represents middleware to wrap handler
type Middleware func(Handler) Handler

func (t *TaskReceiver) taskHandler(ctx context.Context, d amqp.Delivery) error {
	slog.InfoContext(ctx, "start handling")
	var taskDto pb.TaskExecutionMessage

	err := proto.Unmarshal(d.Body, &taskDto)
	if err != nil {
		return err
	}

	translator, err := domain.ParseTranslator(taskDto.Translator.String())
	if err != nil {
		return err
	}
	id, err := uuid.Parse(taskDto.TaskId)
	if err != nil {
		return err
	}

	task := domain.Task{Id: id, Code: taskDto.Code, Translator: translator}

	return t.taskService.Process(ctx, task)
}

func (t *TaskReceiver) applyMiddlewares(h Handler) Handler {
	collapsed := h
	for _, m := range t.middlewares {
		collapsed = m(collapsed)
	}
	return collapsed
}

func (r *TaskReceiver) Receive(ctx context.Context) error {
	msgs, err := r.channel.Consume(
		r.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	h := r.applyMiddlewares(r.taskHandler)

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				return errors.New("rabbitmq connection closed")
			}

			go func(d amqp.Delivery) {
				// Run your handler
				err := h(ctx, d)

				if err != nil {
					slog.ErrorContext(ctx, "handler error", "error", err)
					_ = d.Ack(false)
					return
				}

				if err := d.Ack(false); err != nil {
					slog.Error("failed to ack message", "error", err)
				}
			}(d)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
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
