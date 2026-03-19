package mapper

import (
	"github.com/google/uuid"
	"github.com/kirmala/code_runner/consumer/internal/domain"
	"github.com/kirmala/code_runner/contracts/gen/pb"
)

func ToProtoTask(t domain.Task) pb.TaskExecutionMessage {
	return pb.TaskExecutionMessage{
		TaskId:     t.Id.String(),
		Code:       t.Code,
		Translator: pb.TaskTranslator(t.Translator),
	}
}

func ToDomainTask(t *pb.TaskExecutionMessage) (domain.Task, error) {
	id, err := uuid.Parse(t.TaskId)
	if err != nil {
		return domain.Task{}, ErrInvalidTaskMessage{Field: "id", Err: err.Error()}
	}
	translator := ToDomainTranslator(t.Translator)
	if translator == domain.UnknownTranslator {
		return domain.Task{}, ErrInvalidTaskMessage{Field: "translator", Err: "Unknown translator"}
	}

	return domain.Task{
		Id:   id,
		Code: t.Code,
		Translator: translator,
	}, nil
}
