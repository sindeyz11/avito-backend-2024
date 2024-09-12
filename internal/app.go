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

	//dbConf := config.NewConfig().PostgresConfig()
	//repositories := persistence.NewRepositories(config.NewPostgresConn(dbConf))
	//
	//tenderService := service.NewTenderService(repositories.TenderRepo, repositories.EmployeeRepo)
	//tenderHandler := handlers.NewTenderHandler(tenderService)
	//
	//bidService := service.NewBidService(repositories.BidRepo, repositories.TenderRepo, repositories.EmployeeRepo)
	//bidHandler := handlers.NewBidHandler(bidService)

	mux.HandleFunc("GET /api/ping", handlers.Ping)

	//// tenders
	//mux.HandleFunc("POST /api/tenders/new", tenderHandler.CreateTender)
	//mux.HandleFunc("GET /api/tenders", tenderHandler.GetAllTenders)
	//mux.HandleFunc("GET /api/tenders/my", tenderHandler.GetAllTendersByUsername)
	//mux.HandleFunc("GET /api/tenders/{tenderId}/status", tenderHandler.GetTenderStatusById)
	//mux.HandleFunc("PUT /api/tenders/{tenderId}/status", tenderHandler.UpdateTenderStatusById)
	//mux.HandleFunc("PATCH /api/tenders/{tenderId}/edit", tenderHandler.EditTender)
	//mux.HandleFunc("PUT /api/tenders/{tenderId}/rollback/{version}", tenderHandler.RollbackTender)
	//
	//// bids
	//mux.HandleFunc("POST /api/bids/new", bidHandler.CreateBid)
	//mux.HandleFunc("GET /api/bids/my", bidHandler.GetAllBidsByUsername)
	//
	//mux.HandleFunc("GET /api/bids/{tenderId}/list", bidHandler.GetAllBidsByUsername)
	//mux.HandleFunc("GET /api/bids/{bidId}/status", handlers.Ping)
	//mux.HandleFunc("PUT /api/bids/{bidId}/status", handlers.Ping)
	//mux.HandleFunc("PATCH /api/bids/{bidId}/edit", handlers.Ping)
	//mux.HandleFunc("PUT /api/bids/{bidId}/submit_decision", handlers.Ping)
	//mux.HandleFunc("PUT /api/bids/{bidId}/rollback/{version}", handlers.Ping)

	//serverAddress := os.Getenv("SERVER_ADDRESS")
	//if serverAddress == "" {
	//	serverAddress = ":8080"
	//}
	//fmt.Printf(fmt.Sprintf("Starting server on http://%s/\n", serverAddress))
	fmt.Printf("Starting server on http://127.0.0.1:8080/\n")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
