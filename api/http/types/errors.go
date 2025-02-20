package types

import (
	"net/http"
	"photo_editor/repository"
	"encoding/json"
)

func ProcessError(w http.ResponseWriter, err error, resp any) {
	if err == repository.NotFound {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	} else if err == repository.AlreadyExists {
		http.Error(w, "Key already exitst", http.StatusAlreadyReported)
		return
	}  else if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if resp != nil {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
	}
}