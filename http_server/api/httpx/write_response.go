package httpx

import (
	"code_processor/http_server/api"
	"code_processor/http_server/repository"
	"encoding/json"
	"errors"
	"net/http"
)

type HTTPError struct {
	Error string `json:"error"`
}

func writeErrorMsg(w http.ResponseWriter, msg HTTPError) {
	jsonMsg, _ := json.Marshal(msg)
	_, _ = w.Write(jsonMsg)
}

func WriteResponse(w http.ResponseWriter, err error, response any) {
	w.Header().Set("Content-Type", "application/json")

	if err == nil {
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(response)
		_, _ = w.Write(jsonResponse)
		return
	}

	var (
		notFound repository.ErrNotFound
		conflict repository.ErrConflict
		badReq   api.ErrBadRequest
	)

	switch {
	case errors.As(err, &notFound):
		w.WriteHeader(http.StatusNotFound)
		msg := HTTPError{Error: notFound.Error()}
		writeErrorMsg(w, msg)
	case errors.As(err, &conflict):
		w.WriteHeader(http.StatusConflict)
		msg := HTTPError{Error: conflict.Error()}
		writeErrorMsg(w, msg)
	case errors.As(err, &badReq):
		w.WriteHeader(http.StatusBadRequest)
		msg := HTTPError{Error: badReq.Error()}
		writeErrorMsg(w, msg)
	case errors.Is(err, api.ErrUnauthorized):
		w.WriteHeader(http.StatusUnauthorized)
		msg := HTTPError{Error: "Unauthorized"}
		writeErrorMsg(w, msg)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		msg := HTTPError{Error: "Internal Server Error"}
		writeErrorMsg(w, msg)
	}
}
