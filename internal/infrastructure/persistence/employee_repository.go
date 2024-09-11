package persistence

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"tenders/internal/domain/repository"
	"tenders/internal/utils"
	"tenders/internal/utils/consts"
)

type EmployeeRepo struct {
	Conn *sql.DB
}

func NewEmployeeRepository(conn *sql.DB) *EmployeeRepo {
	return &EmployeeRepo{Conn: conn}
}

var _ repository.EmployeeRepository = &EmployeeRepo{}

func (r *EmployeeRepo) FindEmployeeIdByUsernameIfResponsibleForOrg(username string, organizationId uuid.UUID) (uuid.UUID, error) {
	var employeeId uuid.UUID
	query := `
        SELECT e.id 
        FROM employee e
        JOIN organization_responsible org ON e.id = org.user_id
        WHERE e.username = $1 AND org.organization_id = $2
    `
	err := r.Conn.QueryRow(query, username, organizationId).Scan(&employeeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, errors.New(consts.CannotFindUserError)
		}
		return uuid.Nil, errors.New(consts.UnknownBDError)
	}
	return employeeId, nil
}

func (r *EmployeeRepo) FindEmployeeIdByUsername(username string) (uuid.UUID, error) {
	var employeeId uuid.UUID
	query := `
        SELECT e.id 
        FROM employee e
        WHERE e.username = $1
    `
	err := r.Conn.QueryRow(query, username).Scan(&employeeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, utils.ErrorElementNotExist
		}
		return uuid.Nil, errors.New(consts.UnknownBDError)
	}
	return employeeId, nil
}
