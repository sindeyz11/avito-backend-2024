package interfaces

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
)

type ReviewService interface {
	SubmitFeedback(bidId uuid.UUID, username string, feedback string) (*entity.Bid, error)
	FindAllReviewsByBidAuthor(tenderId uuid.UUID, authorUsername string, requesterUsername string, limit, offset int) ([]entity.Review, error)
}
