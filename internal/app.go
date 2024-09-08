package internal

import (
	"fmt"
	"log"
	"net/http"
	"tenders/internal/interfaces/handlers"
	"tenders/internal/interfaces/middleware"
)

func Run() {
	mux := http.NewServeMux()

	handler := middleware.Logging(mux)

	mux.HandleFunc("/ping", handlers.Ping)

	fmt.Printf("Starting server at port 8080\nhttp://127.0.0.1:8080/\n")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
