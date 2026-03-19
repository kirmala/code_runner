package api

import (
	"fmt"
)

type ErrBadRequest struct {
	Field string
	Err string
}

func (e ErrBadRequest) Error() string {
	if e.Field == "" || e.Err == "" {
		return "bad request"
	}
	msg := fmt.Sprintf("%s: %s", e.Field, e.Err)
	return msg
}

