package request

import (
	"tenders/internal/domain/entity"
	"tenders/internal/utils"
)

type EditBidRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (request EditBidRequest) MapToBid(bid entity.Bid) (entity.Bid, error) {
	var errorFields []string

	if request.Name != "" {
		if len(request.Name) > 100 {
			errorFields = append(errorFields, "name")
		} else {
			bid.Name = request.Name
		}
	}

	if request.Description != "" {
		if len(request.Description) > 500 {
			errorFields = append(errorFields, "description")
		} else {
			bid.Description = request.Description
		}
	}

	if len(errorFields) > 0 {
		return bid, utils.NewValidationError(errorFields)
	}

	return bid, nil
}
