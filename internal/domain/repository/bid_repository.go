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
	FindByTenderIdAndVersion(bidId uuid.UUID, version int) (*entity.Bid, error)
	FindLatestVersionByTenderId(bidId uuid.UUID) (int, error)
	FindAllByOrganizationForEmployee(employeeId uuid.UUID) ([]entity.Bid, error)
	SaveHistoricalVersionTx(tx *sql.Tx, bid *entity.Bid) error
	UpdateStatusAndVersionTx(tx *sql.Tx, bidId uuid.UUID, status string, version int) error
}
