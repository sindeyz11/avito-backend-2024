package persistence

import (
	"database/sql"
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/domain/repository"
	"tenders/internal/utils/consts"
)

type BidRepo struct {
	Conn *sql.DB
}

func NewBidRepository(conn *sql.DB) *BidRepo {
	return &BidRepo{Conn: conn}
}

var _ repository.BidRepository = &BidRepo{}

func (r *BidRepo) Create(bid *entity.Bid) (*entity.Bid, error) {
	query := `INSERT INTO bid (
            bid_id, name, description, status, tender_id,
        	tender_version, author_type, author_id, version, created_at
        ) 
	    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING created_at`
	err := r.Conn.QueryRow(query,
		bid.BidId, bid.Name, bid.Description, bid.Status,
		bid.TenderId, bid.TenderVersion, bid.AuthorType,
		bid.AuthorId, bid.Version, bid.CreatedAt.ConvertToTime(),
	).Scan(&bid.CreatedAt)
	if err != nil {
		return nil, err
	}
	return bid, nil
}

func (r *BidRepo) FindAllByEmployeeIdAndOrgId(employeeId, orgId uuid.UUID, limit, offset int) ([]entity.Bid, error) {
	bids := []entity.Bid{}
	queryStr := `
		SELECT b.bid_id, b.name, b.description, b.status,
		       b.tender_id, b.tender_version, b.author_type, b.author_id, b.version, b.created_at
		FROM bid b
		WHERE (author_id = $1 AND b.author_type = $2) OR (author_id = $3 AND b.author_type = $4)
		ORDER BY name ASC LIMIT $5 OFFSET $6
	`

	rows, err := r.Conn.Query(queryStr,
		employeeId, consts.AuthorTypeUser,
		orgId, consts.AuthorTypeOrganization,
		limit, offset,
	)
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

func (r *BidRepo) FindAllByTenderId(tenderId uuid.UUID, limit, offset int) ([]entity.Bid, error) {
	bids := []entity.Bid{}
	queryStr := `
		SELECT bid_id, name, description, tender_id, tender_version, status, author_type, author_id, version, created_at
		FROM bid
		WHERE tender_id = $1
		ORDER BY name ASC LIMIT $2 OFFSET $3
	`

	rows, err := r.Conn.Query(queryStr, tenderId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bid entity.Bid
		err = rows.Scan(
			&bid.BidId, &bid.Name, &bid.Description,
			&bid.TenderId, &bid.TenderVersion, &bid.Status,
			&bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt,
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

func (r *BidRepo) FindByBidId(bidId uuid.UUID) (*entity.Bid, error) {
	var bid entity.Bid
	queryStr := `
		SELECT bid_id, name, description, tender_id, tender_version, status, author_type, author_id, version, created_at
		FROM bid
		WHERE bid_id = $1
	`

	row := r.Conn.QueryRow(queryStr, bidId)

	err := row.Scan(
		&bid.BidId, &bid.Name, &bid.Description,
		&bid.TenderId, &bid.TenderVersion, &bid.Status,
		&bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &bid, nil
}

func (r *BidRepo) FindByBidIdTx(tx *sql.Tx, bidId uuid.UUID) (*entity.Bid, error) {
	query := `
		SELECT bid_id, name, description, status, tender_id, tender_version, author_type, author_id, version, created_at
		FROM bid
		WHERE bid_id = $1
	`
	row := tx.QueryRow(query, bidId)
	var bid entity.Bid
	err := row.Scan(
		&bid.BidId, &bid.Name, &bid.Description,
		&bid.Status, &bid.TenderId, &bid.TenderVersion,
		&bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

func (r *BidRepo) FindAllByOrganizationForEmployee(employeeId uuid.UUID) ([]entity.Bid, error) {
	query := `
		SELECT b.bid_id, b.name, b.description, b.status, b.tender_id,
		       b.tender_version, b.author_type, b.author_id, b.version, b.created_at
		FROM bid b
		JOIN organization_responsible org ON b.author_id = org.organization_id
		WHERE org.user_id = $1 AND b.author_type = 'ORGANIZATION'
	`

	rows, err := r.Conn.Query(query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []entity.Bid
	for rows.Next() {
		var bid entity.Bid
		if err = rows.Scan(
			&bid.BidId, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId, &bid.TenderVersion,
			&bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt,
		); err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bids, nil
}

func (r *BidRepo) BeginTransaction() (*sql.Tx, error) {
	return r.Conn.Begin()
}

func (r *BidRepo) SaveHistoricalVersionTx(tx *sql.Tx, bid *entity.Bid) error {
	query := `
		INSERT INTO bid_history (
		    bid_id, name, description, status,
		    tender_id, tender_version, author_type, author_id, version, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := tx.Exec(query,
		bid.BidId, bid.Name, bid.Description, bid.Status, bid.TenderId,
		bid.TenderVersion, bid.AuthorType, bid.AuthorId, bid.Version, bid.CreatedAt.ConvertToTime(),
	)
	return err
}

func (r *BidRepo) UpdateBidTx(tx *sql.Tx, bid *entity.Bid) error {
	query := `
		UPDATE bid
		SET name = $1, description = $2, status = $3, tender_id = $4, 
		    tender_version = $5, author_type = $6, author_id = $7, version = $8
		WHERE bid_id = $9
	`

	_, err := tx.Exec(query,
		bid.Name, bid.Description, bid.Status, bid.TenderId,
		bid.TenderVersion, bid.AuthorType, bid.AuthorId, bid.Version,
		bid.BidId,
	)
	return err
}

func (r *BidRepo) FindVersionInHistoryTx(tx *sql.Tx, bidId uuid.UUID, version int) (*entity.Bid, error) {
	query := `
		SELECT 
		    bid_id, name, description, status, tender_id,
		    tender_version, author_type, author_id, version, created_at
		FROM bid_history
		WHERE bid_id = $1 and version = $2
		LIMIT 1
	`
	row := tx.QueryRow(query, bidId, version)
	var bid entity.Bid
	err := row.Scan(
		&bid.BidId, &bid.Name, &bid.Description,
		&bid.Status, &bid.TenderId, &bid.TenderVersion,
		&bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

func (r *BidRepo) FindByAuthorAndTender(authorId uuid.UUID, tenderId uuid.UUID) (*entity.Bid, error) {
	query := `
        SELECT bid_id, name, description, status, tender_id,
               tender_version, author_type, author_id, version, created_at
        FROM bid
        WHERE author_id = $1 AND tender_id = $2
    `
	var bid entity.Bid
	err := r.Conn.QueryRow(query, authorId, tenderId).Scan(
		&bid.BidId, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId,
		&bid.TenderVersion, &bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &bid, nil
}
