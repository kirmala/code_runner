package domain

import "errors"

type Translator int

const (
	UnknownTranslator Translator = iota
	PythonTranslator
	ClangTranslator
	GppTranslator
)

var ErrUnknownTranslator = errors.New("unknown translator")

var TranslatorName = map[Translator]string{
	UnknownTranslator: "unknown",
	PythonTranslator:  "python",
	ClangTranslator:   "clang",
	GppTranslator:     "g++",
}

func ParseTranslator(translator string) (Translator, error) {
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

func (t Translator) String() string {
	return TranslatorName[t]
}