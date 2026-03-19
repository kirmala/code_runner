package mapper

import (
	"fmt"
)

type ErrInvalidTaskMessage struct {
	Field string
	Err string
}

func (e ErrInvalidTaskMessage) Error() string {
	if e.Field == "" || e.Err == "" {
		return "bad request"
	}
	msg := fmt.Sprintf("%s: %s", e.Field, e.Err)
	return msg
}