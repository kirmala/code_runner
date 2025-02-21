package editing

import (
	"time"
	"photo_editor/usecases"
	"photo_editor/models"
)

func Edit(service usecases.Task, newTask *models.Task) {
	time.Sleep(8 * time.Second)
	newTask.Status = "ready"
	newTask.Result = "something happend"
	service.Put(*newTask)
}