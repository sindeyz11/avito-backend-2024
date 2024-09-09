package repository

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
)

type TenderRepository interface {
	Create(tender *entity.Tender) (*entity.Tender, error)
	FindByEmployeeUsername(username string) (*entity.Tender, error)
	FindAll() ([]*entity.Tender, error)
	Update(tender *entity.Tender) (*entity.Tender, error)
	Delete(id uuid.UUID) error
}
