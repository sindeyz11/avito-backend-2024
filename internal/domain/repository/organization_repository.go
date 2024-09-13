package repository

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
)

type OrganizationRepository interface {
	FindById(id uuid.UUID) (*entity.Organization, error)
	FindByEmployeeId(employeeId uuid.UUID) (*entity.Organization, error)
}
