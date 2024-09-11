package handlers

import (
	"fmt"
	"net/http"
	"tenders/internal/utils/consts"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprintln(w, "ok")
	if err != nil {
		http.Error(w, consts.FailedToWriteResponse, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
