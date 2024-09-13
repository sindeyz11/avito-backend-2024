package request

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/utils"
	"tenders/internal/utils/consts"
)

type TenderRequest struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ServiceType     string    `json:"serviceType"`
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

	if tenderRequest.ServiceType != consts.Construction &&
		tenderRequest.ServiceType != consts.Delivery &&
		tenderRequest.ServiceType != consts.Manufacture {
		errorFields = append(errorFields, "serviceType")
	}

	if tenderRequest.OrganizationID.String() == "" {
		errorFields = append(errorFields, "organizationID")
	}

	if len(errorFields) > 0 {
		return nil, utils.NewValidationError(errorFields)
	}

	tender := entity.Tender{
		Name:           tenderRequest.Name,
		Description:    tenderRequest.Description,
		ServiceType:    tenderRequest.ServiceType,
		OrganizationID: tenderRequest.OrganizationID,
	}
	return &tender, nil
}
