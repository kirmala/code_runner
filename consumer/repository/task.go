package repository

import "consumer/models"

type Task interface {
	Put(models.Task) error
}