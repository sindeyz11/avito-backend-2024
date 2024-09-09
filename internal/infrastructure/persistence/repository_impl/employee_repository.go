package repository_impl

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"tenders/internal/domain/repository"
	"tenders/internal/utils/consts"
)

type EmployeeRepo struct {
	Conn *sql.DB
}

func NewEmployeeRepository(conn *sql.DB) *EmployeeRepo {
	return &EmployeeRepo{Conn: conn}
}

var _ repository.EmployeeRepository = &EmployeeRepo{}

func (r *EmployeeRepo) GetEmployeeIdByUsernameIfResponsibleForOrg(username string, organizationId uuid.UUID) (uuid.UUID, error) {
	var creatorID uuid.UUID
	query := `
        SELECT e.id 
        FROM employee e
        JOIN organization_responsible org ON e.id = org.user_id
        WHERE e.username = $1 AND org.organization_id = $2
    `
	err := r.Conn.QueryRow(query, username, organizationId).Scan(&creatorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, errors.New(consts.UserNotExistsError)
		}
		return uuid.Nil, errors.New("database query error")
	}
	return creatorID, nil
}
