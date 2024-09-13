package request

import (
	"tenders/internal/domain/entity"
	"tenders/internal/utils"
	"tenders/internal/utils/consts"
)

type EditTenderRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"serviceType"`
}

func (request EditTenderRequest) MapToTender(tender *entity.Tender) error {
	var errorFields []string

	if request.Name != "" {
		if len(request.Name) > 100 {
			errorFields = append(errorFields, "name")
		} else {
			tender.Name = request.Name
		}
	}

	if request.Description != "" {
		if len(request.Description) > 500 {
			errorFields = append(errorFields, "description")
		} else {
			tender.Description = request.Description
		}
	}

	switch request.ServiceType {
	case "":
		// Если не передан, оставляем текущее значение
	case consts.Construction, consts.Delivery, consts.Manufacture:
		tender.ServiceType = request.ServiceType
	default:
		errorFields = append(errorFields, "serviceType")
	}

	if len(errorFields) > 0 {
		return utils.NewValidationError(errorFields)
	}

	return nil
}
