package repository

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
)

type EmployeeRepository interface {
	FindEmployeeIdByUsernameIfResponsibleForOrg(username string, organizationId uuid.UUID) (uuid.UUID, error)
	FindEmployeeIdByUsername(username string) (uuid.UUID, error)
	FindById(id uuid.UUID) (*entity.Employee, error)
}
