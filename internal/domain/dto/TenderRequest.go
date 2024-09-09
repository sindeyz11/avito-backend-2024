package dto

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/utils"
)

type TenderRequest struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ServiceType     string    `json:"serviceType"`
	Status          string    `json:"status"`
	OrganizationID  uuid.UUID `json:"organizationId"`
	CreatorUsername string    `json:"creatorUsername"`
}

func (tenderRequest TenderRequest) MapToTender() (*entity.Tender, error) {
	var errorFields []string

	if len(tenderRequest.Description) > 500 {
		errorFields = append(errorFields, "description")
	}

	if tenderRequest.ServiceType != entity.CONSTRUCTION &&
		tenderRequest.ServiceType != entity.DELIVERY &&
		tenderRequest.ServiceType != entity.MANUFACTURE {
		errorFields = append(errorFields, "serviceType")
	}

	if len(errorFields) > 0 {
		return nil, utils.NewValidationError(errorFields)
	}

	tender := entity.Tender{
		Name:           tenderRequest.Name,
		Description:    tenderRequest.Description,
		ServiceType:    tenderRequest.ServiceType,
		Status:         tenderRequest.Status,
		OrganizationID: tenderRequest.OrganizationID,
	}
	return &tender, nil
}
