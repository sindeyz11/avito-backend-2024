package interfaces

import (
	"tenders/internal/domain/dto"
	"tenders/internal/domain/entity"
)

type TenderService interface {
	CreateTender(tenderRequest *dto.TenderRequest) (*entity.Tender, error)
}
