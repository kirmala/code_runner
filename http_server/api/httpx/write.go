package httpx

import (
	"code_processor/http_server/api"
	"code_processor/http_server/repository"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err == nil {
		return
	}

	var (
		notFound repository.ErrNotFound
		conflict repository.ErrConflict
		badReq   api.ErrBadRequest
	)

	switch {
	case errors.As(err, &notFound):
		msg := HTTPError{Error: notFound.Error()}
		writeError(w, msg, http.StatusNotFound)
	case errors.As(err, &conflict):
		msg := HTTPError{Error: conflict.Error()}
		writeError(w, msg, http.StatusConflict)
	case errors.As(err, &badReq):
		msg := HTTPError{Error: badReq.Error()}
		writeError(w, msg, http.StatusBadRequest)
	case errors.Is(err, api.ErrUnauthorized):
		msg := HTTPError{Error: "Unauthorized"}
		writeError(w, msg, http.StatusUnauthorized)
	default:
		msg := HTTPError{Error: "Internal Server Error"}
		writeError(w, msg, http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, msg HTTPError, errorStatus int) {
	if errorStatus < 400 || errorStatus >= 600 {
		panic(fmt.Sprintf("invalid error status: %d", errorStatus))
	}

	w.WriteHeader(int(errorStatus))
	jsonMsg, _ := json.Marshal(msg)
	_, _ = w.Write(jsonMsg)
}

func WriteSuccess(w http.ResponseWriter, response any, successStatus int) {
	if successStatus < 200 || successStatus >= 300 {
		panic(fmt.Sprintf("invalid success status: %d", successStatus))
	}

	if response != nil {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(int(successStatus))

	if response != nil {
		jsonResponse, _ := json.Marshal(response)
		_, _ = w.Write(jsonResponse)
	}
}
