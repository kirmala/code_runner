package types

import (
	"code_processor/http_server/repository"
	"encoding/json"
	"net/http"
)

func ProcessError(w http.ResponseWriter, err error, resp any, successResponse int) {
	if err == repository.ErrNotFound {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	} else if err == repository.ErrAlreadyExists {
		http.Error(w, "Key already exists", http.StatusAlreadyReported)
		return
	} else if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(successResponse)
	if resp != nil {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	}
}
