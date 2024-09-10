package internal

import (
	"fmt"
	"log"
	"net/http"
	"tenders/internal/application/service"
	"tenders/internal/infrastructure/config"
	"tenders/internal/infrastructure/persistence"
	"tenders/internal/interfaces/handlers"
	"tenders/internal/interfaces/middleware"
)

func Run() {
	mux := http.NewServeMux()
	handler := middleware.Logging(mux)

	dbConf := config.NewConfig().PostgresConfig()
	repositories := persistence.NewRepositories(config.NewPostgresConn(dbConf))

	tenderService := service.NewTenderService(repositories.TenderRepo, repositories.EmployeeRepo)
	tenderController := handlers.NewTenderHandler(tenderService)

	mux.HandleFunc("/api/ping", handlers.Ping)
	mux.HandleFunc("POST /api/tenders/new", tenderController.CreateTender)
	mux.HandleFunc("GET /api/tenders", tenderController.GetAllTenders)
	mux.HandleFunc("GET /api/tenders/my", tenderController.GetAllTendersByUsername)
	mux.HandleFunc("GET /api/tenders/{tenderId}/status", tenderController.GetTenderStatusById)
	mux.HandleFunc("PUT /api/tenders/{tenderId}/status", tenderController.UpdateTenderStatusById)

	fmt.Printf("Starting server at port 8080\nhttp://127.0.0.1:8080/\n")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
