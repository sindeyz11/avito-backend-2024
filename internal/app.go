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

	tenderService := service.NewTenderService(repositories.TenderRepo, repositories.EmployeeRepo, repositories.OrganizationRepo)
	tenderHandler := handlers.NewTenderHandler(tenderService)

	bidService := service.NewBidService(
		repositories.EmployeeRepo, repositories.OrganizationRepo, repositories.BidRepo, repositories.TenderRepo,
	)
	bidHandler := handlers.NewBidHandler(bidService)

	reviewService := service.NewReviewService(
		repositories.EmployeeRepo, repositories.OrganizationRepo,
		repositories.BidRepo, repositories.TenderRepo, repositories.ReviewRepo,
	)
	feedbackHandler := handlers.NewReviewHandler(reviewService)

	mux.HandleFunc("GET /api/ping", handlers.Ping)

	// tenders
	mux.HandleFunc("POST /api/tenders/new", tenderHandler.CreateTender)
	mux.HandleFunc("GET /api/tenders", tenderHandler.GetAllTenders)
	mux.HandleFunc("GET /api/tenders/my", tenderHandler.GetAllTendersByUsername)
	mux.HandleFunc("GET /api/tenders/{tenderId}/status", tenderHandler.GetTenderStatusById)
	mux.HandleFunc("PUT /api/tenders/{tenderId}/status", tenderHandler.UpdateTenderStatusById)
	mux.HandleFunc("PATCH /api/tenders/{tenderId}/edit", tenderHandler.EditTender)
	mux.HandleFunc("PUT /api/tenders/{tenderId}/rollback/{version}", tenderHandler.RollbackTender)

	// bids
	mux.HandleFunc("POST /api/bids/new", bidHandler.CreateBid)
	mux.HandleFunc("GET /api/bids/my", bidHandler.GetAllBidsByUsername)
	mux.HandleFunc("GET /api/bids/{tenderId}/list", bidHandler.GetAllBidsByTender)
	mux.HandleFunc("GET /api/bids/{bidId}/status", bidHandler.GetBidStatusById)
	mux.HandleFunc("PUT /api/bids/{bidId}/status", bidHandler.UpdateBidStatusById)
	mux.HandleFunc("PATCH /api/bids/{bidId}/edit", bidHandler.EditBid)
	mux.HandleFunc("PUT /api/bids/{bidId}/rollback/{version}", bidHandler.RollbackBid)
	mux.HandleFunc("PUT /api/bids/{bidId}/submit_decision", bidHandler.SubmitDecision)

	// feedback
	mux.HandleFunc("GET /api/bids/{tenderId}/reviews", feedbackHandler.GetReviewsList)
	mux.HandleFunc("PUT /api/bids/{bidId}/feedback", feedbackHandler.SubmitFeedback)

	fmt.Printf("Starting server on http://127.0.0.1:8080/\n")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
