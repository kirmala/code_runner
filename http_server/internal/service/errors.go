package service

import "fmt"

type ErrUnauthenticated struct {
	Msg string
}

func (e ErrUnauthenticated) Error() string {
	if e.Msg == "" {
		return "unauthenticated"
	}
	msg := fmt.Sprintf("unauthenticated: %s", e.Msg)
	return msg
}