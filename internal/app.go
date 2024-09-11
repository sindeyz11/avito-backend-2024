package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

	mux.HandleFunc("GET /api/ping", handlers.Ping)

	// tenders
	mux.HandleFunc("POST /api/tenders/new", tenderController.CreateTender)
	mux.HandleFunc("GET /api/tenders", tenderController.GetAllTenders)
	mux.HandleFunc("GET /api/tenders/my", tenderController.GetAllTendersByUsername)
	mux.HandleFunc("GET /api/tenders/{tenderId}/status", tenderController.GetTenderStatusById)
	mux.HandleFunc("PUT /api/tenders/{tenderId}/status", tenderController.UpdateTenderStatusById)
	mux.HandleFunc("PATCH /api/tenders/{tenderId}/edit", tenderController.EditTender)
	mux.HandleFunc("PUT /api/tenders/{tenderId}/rollback/{version}", tenderController.RollbackTender)

	// bids
	mux.HandleFunc("GET /api/bids", handlers.Ping)

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		log.Fatalf("Can't get server address from env")
	}
	fmt.Printf(fmt.Sprintf("Starting server on http://%s/\n", serverAddress))

	if err := http.ListenAndServe(serverAddress, handler); err != nil {
		log.Fatal(err)
	}
}
