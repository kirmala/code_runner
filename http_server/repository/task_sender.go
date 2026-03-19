package repository

import "github.com/kirmala/code_runner/http_server/domain"

type TaskSender interface {
	Send(domain.Task) error
}
