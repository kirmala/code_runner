package domain

import (
	"fmt"
)

type Translator int

const (
	UnknownTranslator Translator = iota
	PythonTranslator
	ClangTranslator
	GppTranslator
)

type ErrUnknownTranslator struct {
	translator string
}

func (e ErrUnknownTranslator) Error() string {
	return fmt.Sprintf("unknown translator: %s", e.translator)
}

var TranslatorName = map[Translator]string{
	UnknownTranslator: "unknown",
	PythonTranslator:  "python",
	ClangTranslator:   "clang",
	GppTranslator:     "g++",
}

func ParseTranslator(translator string) (Translator, error) {
	switch translator {
	case "TASK_TRANSLATOR_PYTHON":
		return PythonTranslator, nil
	case "TASK_TRANSLATOR_CLANG":
		return ClangTranslator, nil
	case "TASK_TRANSLATOR_GPP":
		return GppTranslator, nil
	default:
		return UnknownTranslator, ErrUnknownTranslator{translator: translator}
	}
}

func (t Translator) String() string {
	return TranslatorName[t]
}