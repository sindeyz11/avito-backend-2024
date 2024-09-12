package interfaces

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/interfaces/dto/request"
)

type BidService interface {
	FindAllByEmployeeUsername(username string, limit, offset int) ([]entity.Bid, error)
	FindAllByTenderId(tenderId uuid.UUID, username string, limit, offset int) ([]entity.Bid, error)
	Create(request *request.BidRequest) (*entity.Bid, error)
	GetStatusByBidId(bidId uuid.UUID, username string) (string, error)
}
