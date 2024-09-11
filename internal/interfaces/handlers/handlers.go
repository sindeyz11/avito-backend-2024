package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tenders/internal/domain/dto"
)

func HandleError(w http.ResponseWriter, statusCode int, errorMsg string) {
	j, err := json.Marshal(dto.ErrorResponse{Reason: errorMsg})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	_, _ = w.Write(j)
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "ok")
	if err != nil {
		return
	}
}
