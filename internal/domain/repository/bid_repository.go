package repository

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
)

type BidRepository interface {
	Create(tender *entity.Bid) (*entity.Bid, error)
	FindAllByEmployeeId(id uuid.UUID, limit, offset int) ([]entity.Bid, error)
	FindAllByTenderId(tenderId uuid.UUID, limit, offset int) ([]entity.Bid, error)
	FindByBidId(bidId uuid.UUID) (*entity.Bid, error)
	FindByTenderIdAndVersion(bidId uuid.UUID, version int) (*entity.Bid, error)
	FindLatestVersionByTenderId(bidId uuid.UUID) (int, error)
}
