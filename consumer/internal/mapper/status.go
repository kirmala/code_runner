package mapper

import (
	"github.com/kirmala/code_runner/consumer/internal/domain"
	"github.com/kirmala/code_runner/contracts/gen/pb"
)

func ToProtoStatus(t domain.Status) pb.TaskStatus {
    switch t {
    case domain.StatusCompleted:
        return pb.TaskStatus_TASK_STATUS_COMPLETED
    case domain.StatusInProgress:
        return pb.TaskStatus_TASK_STATUS_IN_PROGRESS
    case domain.StatusFailed:
        return pb.TaskStatus_TASK_STATUS_FAILED
    default:
        return pb.TaskStatus_TASK_STATUS_UNKNOWN
    }
}

func ToDomainStatus(t pb.TaskStatus) domain.Status {
    switch t {
    case pb.TaskStatus_TASK_STATUS_COMPLETED:
        return domain.StatusCompleted
    case pb.TaskStatus_TASK_STATUS_IN_PROGRESS:
        return domain.StatusInProgress
    case pb.TaskStatus_TASK_STATUS_FAILED:
        return domain.StatusFailed
    default:
        return domain.UnknownStatus
    }
}