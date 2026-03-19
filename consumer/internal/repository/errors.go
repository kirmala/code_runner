package repository

import "fmt"

type ErrNotFound struct {
	Item string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found", e.Item)
}