package models

import (
	"errors"
)

type status int

const (
	UnknownStatus status = iota
	StatusInProgress
	StatusCompleted
	StatusFailed
)

var ErrUnknownStatus = errors.New("unknown status")

var StatusName = map[status]string{
	UnknownStatus: "unknown",
	StatusInProgress: "in_progress",
	StatusCompleted:  "completed",
	StatusFailed:    "failed",
}

func ParseStatus(status string) (status, error) {
	switch status {
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

func (s status) String() string {
	switch s {
	case StatusInProgress:
		return "in_progress"
	case StatusCompleted:
		return "completed"
	case StatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}



