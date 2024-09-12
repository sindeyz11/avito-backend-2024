package persistence

import (
	"database/sql"
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/domain/repository"
)

type BidRepo struct {
	Conn *sql.DB
}

func NewBidRepository(conn *sql.DB) *BidRepo {
	return &BidRepo{Conn: conn}
}

var _ repository.BidRepository = &BidRepo{}

func (r *BidRepo) Create(bid *entity.Bid) (*entity.Bid, error) {
	query := `INSERT INTO bid (bid_id, name, description, status, tender_id,
                 tender_version, author_type, author_id, version, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.Conn.Exec(query,
		bid.BidId, bid.Name, bid.Description, bid.Status,
		bid.TenderId, bid.TenderVersion, bid.AuthorType,
		bid.AuthorId, bid.Version, bid.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return bid, nil
}

func (r *BidRepo) FindAllByEmployeeId(id uuid.UUID, limit, offset int) ([]entity.Bid, error) {
	bids := []entity.Bid{}
	queryStr := `
		SELECT DISTINCT ON (b.bid_id, b.name) b.bid_id, b.name, b.description, b.status,
		       b.tender_id, b.tender_version, b.author_type, b.author_id, b.version, b.created_at
		FROM bid b
		WHERE b.author_id = $1 AND b.version = (
			SELECT MAX(version)
			FROM bid
			WHERE bid_id = b.bid_id
		)
		ORDER BY name ASC LIMIT $2 OFFSET $3
	`

	rows, err := r.Conn.Query(queryStr, id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bid entity.Bid
		err = rows.Scan(
			&bid.BidId, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId,
			&bid.TenderVersion, &bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		bids = append(bids, bid)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bids, nil
}

func (r *BidRepo) FindAllByTenderId(tenderId uuid.UUID) ([]entity.Bid, error) {
	//TODO implement me
	panic("implement me")
}

func (r *BidRepo) FindByBidId(bidId uuid.UUID) (*entity.Bid, error) {
	//TODO implement me
	panic("implement me")
}

func (r *BidRepo) FindByTenderIdAndVersion(bidId uuid.UUID, version int) (*entity.Bid, error) {
	//TODO implement me
	panic("implement me")
}

func (r *BidRepo) FindLatestVersionByTenderId(bidId uuid.UUID) (int, error) {
	//TODO implement me
	panic("implement me")
}
