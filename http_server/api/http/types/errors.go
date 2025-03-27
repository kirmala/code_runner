package types

import (
	"code_processor/http_server/repository"
	"encoding/json"
	"net/http"
)

func ProcessError(w http.ResponseWriter, err error, resp any, successResponse int) {
	if err == repository.NotFound {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	} else if err == repository.AlreadyExists {
		http.Error(w, "Key already exists", http.StatusAlreadyReported)
		return
	} else if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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
