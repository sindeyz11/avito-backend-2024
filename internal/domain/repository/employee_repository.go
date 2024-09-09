package repository

import "github.com/google/uuid"

type EmployeeRepository interface {
	GetEmployeeIdByUsernameIfResponsibleForOrg(username string, organizationId uuid.UUID) (uuid.UUID, error)
}
