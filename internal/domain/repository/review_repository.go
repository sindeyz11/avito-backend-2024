package repository

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
)

type ReviewRepository interface {
	Create(*entity.Review) (*entity.Review, error)
	FindAllByBidAuthor(authorId uuid.UUID, limit, offset int) ([]entity.Review, error)
}
