package dto

type GetTaskHandlerRequest struct {
	Id string `json:"id"`
}

type GetTaskStatusHandlerResponse struct {
	Status string `json:"status"`
}

type GetTaskResultHandlerResponse struct {
	Result *string `json:"result"`
}

type PostTaskHandlerRequest struct {
	TaskCode       string `json:"task_code"`
	TaskTranslator string `json:"task_translator"`
}

type PostTaskHandlerResponse struct {
	ID string `json:"id"`
}
