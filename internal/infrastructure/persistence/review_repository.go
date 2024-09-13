package persistence

import (
	"database/sql"
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/domain/repository"
)

type ReviewRepo struct {
	Conn *sql.DB
}

func NewReviewRepository(conn *sql.DB) *ReviewRepo {
	return &ReviewRepo{Conn: conn}
}

var _ repository.ReviewRepository = &ReviewRepo{}

func (r *ReviewRepo) Create(review *entity.Review) (*entity.Review, error) {
	query := `INSERT INTO review (id, bid_id, description, created_at) 
	    VALUES ($1, $2, $3, $4) RETURNING created_at`
	err := r.Conn.QueryRow(query,
		review.Id, review.BidId, review.Description, review.CreatedAt.ConvertToTime(),
	).Scan(&review.CreatedAt)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r *ReviewRepo) FindAllByBidAuthor(authorId uuid.UUID, limit, offset int) ([]entity.Review, error) {
	query := `
        SELECT id, bid_id, description, created_at
        FROM review
        WHERE bid_id IN (
            SELECT bid_id
            FROM bid
            WHERE author_id = $1
        )
        ORDER BY description ASC
        LIMIT $2 OFFSET $3
    `
	rows, err := r.Conn.Query(query, authorId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := []entity.Review{}
	for rows.Next() {
		var review entity.Review
		err = rows.Scan(&review.Id, &review.BidId, &review.Description, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
