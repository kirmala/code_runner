package domain

import (
	"errors"
)

type Status int

const (
	UnknownStatus Status = iota
	StatusInProgress
	StatusCompleted
	StatusFailed
)

var ErrUnknownStatus = errors.New("unknown Status")

var StatusName = map[Status]string{
	UnknownStatus:    "unknown",
	StatusInProgress: "in_progress",
	StatusCompleted:  "completed",
	StatusFailed:     "failed",
}

func ParseStatus(Status string) (Status, error) {
	switch Status {
	case "in_progress":
		return StatusInProgress, nil
	case "completed":
		return StatusCompleted, nil
	case "failed":
		return StatusFailed, nil
	default:
		return UnknownStatus, ErrUnknownStatus
	}
}

func (s Status) String() string {
	return StatusName[s]
}
