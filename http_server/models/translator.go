package models

import "errors"

type translator int

const (
	UnknownTranslator translator = iota
	PythonTranslator
	ClangTranslator
	GppTranslator
)

var ErrUnknownTranslator = errors.New("unknown translator")

var TranslatorName = map[translator]string{
	UnknownTranslator: "unknown",
	PythonTranslator: "python",
	ClangTranslator:  "clang",
	GppTranslator:    "g++",
}

func ParseTranslator(translator string) (translator, error) {
	switch translator {
	case "python":
		return PythonTranslator, nil
	case "clang":
		return ClangTranslator, nil
	case "g++":
		return GppTranslator, nil
	default:
		return UnknownTranslator, ErrUnknownTranslator
	}
}

func (t translator) String() string {
	switch t {
	case PythonTranslator:
		return "python"
	case ClangTranslator:
		return "clang"
	case GppTranslator:
		return "g++"
	default:
		return "unknown"
	}
}



