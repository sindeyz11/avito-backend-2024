package interfaces

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/interfaces/dto/request"
)

type BidService interface {
	FindAllByEmployeeUsername(username string, limit, offset int) ([]entity.Bid, error)
	FindAllByTenderId(tenderId uuid.UUID, username string, limit, offset int) ([]entity.Bid, error)
	CreateNewBid(request *request.BidRequest) (*entity.Bid, error)
	GetStatusByBidId(bidId uuid.UUID, username string) (string, error)
	UpdateStatus(bidId uuid.UUID, status string, username string) (*entity.Bid, error)
	EditBid(bidId uuid.UUID, username string, updateRequest *request.EditBidRequest) (*entity.Bid, error)
	RollbackBid(bidId uuid.UUID, version int, username string) (*entity.Bid, error)
	SubmitDecision(bidId uuid.UUID, username, decision string) (*entity.Bid, error)
}
