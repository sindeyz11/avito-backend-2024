package repository

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
)

type TenderRepository interface {
	Create(tender *entity.Tender) (*entity.Tender, error)
	FindAllByEmployeeId(id uuid.UUID, limit, offset int) ([]entity.Tender, error)
	FindAll(serviceTypes []string, limit, offset int) ([]entity.Tender, error)
	FindByTenderId(id uuid.UUID) (*entity.Tender, error)
}
