package request

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/utils"
)

type BidRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TenderId    uuid.UUID `json:"tenderId"`
	AuthorType  string    `json:"authorType"`
	AuthorId    uuid.UUID `json:"authorId"`
}

// MapToBid мапит в тендер и валидирует
func (bidRequest BidRequest) MapToBid() (*entity.Bid, error) {
	var errorFields []string

	if bidRequest.Name == "" || len(bidRequest.Name) > 100 {
		errorFields = append(errorFields, "name")
	}

	if bidRequest.Description == "" || len(bidRequest.Description) > 500 {
		errorFields = append(errorFields, "description")
	}

	if bidRequest.AuthorType != entity.ORGANIZATION && bidRequest.AuthorType != entity.USER {
		errorFields = append(errorFields, "authorType")
	}

	if len(errorFields) > 0 {
		return nil, utils.NewValidationError(errorFields)
	}

	return &entity.Bid{
		Name:        bidRequest.Name,
		Description: bidRequest.Description,
		TenderId:    bidRequest.TenderId,
		AuthorType:  bidRequest.AuthorType,
		AuthorId:    bidRequest.AuthorId,
	}, nil
}
