package repository

import "github.com/google/uuid"

type EmployeeRepository interface {
	FindEmployeeIdByUsernameIfResponsibleForOrg(username string, organizationId uuid.UUID) (uuid.UUID, error)
	FindEmployeeIdByUsername(username string) (uuid.UUID, error)
}
