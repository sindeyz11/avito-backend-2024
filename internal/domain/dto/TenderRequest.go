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

// MapToTender мапит в тендер и валидирует
func (tenderRequest TenderRequest) MapToTender() (*entity.Tender, error) {
	var errorFields []string

	if tenderRequest.Name == "" || len(tenderRequest.Name) > 100 {
		errorFields = append(errorFields, "name")
	}

	if tenderRequest.Description == "" || len(tenderRequest.Description) > 500 {
		errorFields = append(errorFields, "description")
	}

	if tenderRequest.ServiceType != entity.CONSTRUCTION &&
		tenderRequest.ServiceType != entity.DELIVERY &&
		tenderRequest.ServiceType != entity.MANUFACTURE {
		errorFields = append(errorFields, "serviceType")
	}

	if tenderRequest.Status != entity.CREATED &&
		tenderRequest.Status != entity.PUBLISHED &&
		tenderRequest.Status != entity.CLOSED {
		errorFields = append(errorFields, "status")
	}

	//if tenderRequest.OrganizationID.String() == "" {
	//	errorFields = append(errorFields, "organizationID")
	//}

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
