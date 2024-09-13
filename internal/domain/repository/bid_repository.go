package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
)

type BidRepository interface {
	BeginTransaction() (*sql.Tx, error)
	Create(bid *entity.Bid) (*entity.Bid, error)
	FindAllByEmployeeIdAndOrgId(employeeId, orgId uuid.UUID, limit, offset int) ([]entity.Bid, error)
	FindAllByTenderId(tenderId uuid.UUID, limit, offset int) ([]entity.Bid, error)
	FindByBidId(bidId uuid.UUID) (*entity.Bid, error)
	FindAllByOrganizationForEmployee(employeeId uuid.UUID) ([]entity.Bid, error)
	SaveHistoricalVersionTx(tx *sql.Tx, bid *entity.Bid) error
	UpdateStatusAndVersionTx(tx *sql.Tx, bidId uuid.UUID, status string, version int) error
	UpdateBidTx(tx *sql.Tx, bid *entity.Bid) error
	FindByBidIdTx(tx *sql.Tx, bidId uuid.UUID) (*entity.Bid, error)
	FindVersionInHistoryTx(tx *sql.Tx, bidId uuid.UUID, version int) (*entity.Bid, error)
	FindByAuthorAndTender(authorId uuid.UUID, tenderId uuid.UUID) (*entity.Bid, error)
}
