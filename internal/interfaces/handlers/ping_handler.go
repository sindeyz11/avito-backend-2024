package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tenders/internal/utils/consts"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
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
		http.Error(w, consts.FailedToWriteResponse, http.StatusInternalServerError)
		return
	}
}
