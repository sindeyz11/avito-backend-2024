package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	if r.Method != http.MethodGet {
		HandleError(w, 400, "Добавить константу")
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "ok\n"+
		os.Getenv("SERVER_ADDRESS")+"\n"+
		os.Getenv("POSTGRES_CONN")+"\n"+
		os.Getenv("POSTGRES_JDBC_URL")+"\n"+
		os.Getenv("POSTGRES_USERNAME")+"\n"+
		os.Getenv("POSTGRES_PASSWORD")+"\n"+
		os.Getenv("POSTGRES_HOST")+"\n"+
		os.Getenv("POSTGRES_PORT")+"\n"+
		os.Getenv("POSTGRES_DATABASE")+"\n"+
		"values")
	if err != nil {
		// TODO обработка internal
		return
	}
}
