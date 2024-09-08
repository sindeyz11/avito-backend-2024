package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
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
		return
	}
}
