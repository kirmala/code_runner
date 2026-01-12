package byHttp

import (
	"bytes"
	"code_processor/http_server/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Task struct {
	addr string
}

func NewTask(addr string) *Task {
	return &Task{addr: addr}
}

func (r *Task) Put(task models.Task) error {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshalling JSON: %s", err)
	}
	resp, err := http.Post(r.addr, "application/json", bytes.NewBuffer(jsonTask))
	if err != nil {
		return fmt.Errorf("sending request: %s", err)
	}
	defer func() {_ = resp.Body.Close()}()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}
	return nil
}
