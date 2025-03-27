package models

type Task struct {
	Id         string `json:"id"`
	Code       string `json:"code"`
	Translator string `json:"translator"`
	Status     string `json:"status"`
	Result     string `json:"result"`
}
