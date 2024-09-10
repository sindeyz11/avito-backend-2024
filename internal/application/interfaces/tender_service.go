package interfaces

import (
	"github.com/google/uuid"
	"tenders/internal/domain/dto"
	"tenders/internal/domain/entity"
)

type TenderService interface {
	Create(tenderRequest *dto.TenderRequest) (*entity.Tender, error)
	FindAll(serviceTypes []string, limit, offset int) ([]entity.Tender, error)
	FindAllByEmployeeUsername(username string, limit, offset int) ([]entity.Tender, error)
	GetStatusByTenderId(id uuid.UUID, username string) (string, error)
	UpdateStatus(tenderId uuid.UUID, status, username string) (*entity.Tender, error)
}
