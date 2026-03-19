package mapper

import (
	"github.com/kirmala/code_runner/consumer/internal/domain"
	"github.com/kirmala/code_runner/contracts/gen/pb"
)

func ToProtoTranslator(t domain.Translator) pb.TaskTranslator {
    switch t {
    case domain.PythonTranslator:
        return pb.TaskTranslator_TASK_TRANSLATOR_PYTHON
    case domain.ClangTranslator:
        return pb.TaskTranslator_TASK_TRANSLATOR_CLANG
    case domain.GppTranslator:
        return pb.TaskTranslator_TASK_TRANSLATOR_GPP
    default:
        return pb.TaskTranslator_TASK_TRANSLATOR_UNKNOWN
    }
}

func ToDomainTranslator(t pb.TaskTranslator) domain.Translator {
    switch t {
    case pb.TaskTranslator_TASK_TRANSLATOR_PYTHON:
        return domain.PythonTranslator
    case pb.TaskTranslator_TASK_TRANSLATOR_CLANG:
        return domain.ClangTranslator
    case pb.TaskTranslator_TASK_TRANSLATOR_GPP:
        return domain.GppTranslator
    default:
        return domain.UnknownTranslator
    }
}