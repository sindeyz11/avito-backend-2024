package repository

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
)

type TenderRepository interface {
	Create(tender *entity.Tender) (*entity.Tender, error)
	FindAllAvailableByOrganizationId(id uuid.UUID, limit, offset int) ([]entity.Tender, error)
	FindAllPublished(serviceTypes []string, limit, offset int) ([]entity.Tender, error)
	FindByTenderId(id uuid.UUID) (*entity.Tender, error)
	FindByTenderIdAndVersion(tenderId uuid.UUID, version int) (*entity.Tender, error)
	FindLatestVersionByTenderId(tenderId uuid.UUID) (int, error)
	CheckEmployeeAccessToTender(employeeId uuid.UUID, tenderId uuid.UUID) (bool, error)
}
