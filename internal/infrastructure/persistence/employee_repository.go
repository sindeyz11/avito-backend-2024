package persistence

import (
	"database/sql"
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/domain/repository"
)

type EmployeeRepo struct {
	Conn *sql.DB
}

func NewEmployeeRepository(conn *sql.DB) *EmployeeRepo {
	return &EmployeeRepo{Conn: conn}
}

var _ repository.EmployeeRepository = &EmployeeRepo{}

func (r *EmployeeRepo) FindEmployeeIdByUsername(username string) (uuid.UUID, error) {
	var employeeId uuid.UUID
	query := `
        SELECT e.id 
        FROM employee e
        WHERE e.username = $1
    `
	err := r.Conn.QueryRow(query, username).Scan(&employeeId)
	if err != nil {
		return uuid.Nil, err
	}
	return employeeId, nil
}

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
		return uuid.Nil, err
	}
	return employeeId, nil
}

func (r *EmployeeRepo) FindById(id uuid.UUID) (*entity.Employee, error) {
	query := `SELECT id, username, first_name, last_name, created_at, updated_at FROM employee WHERE id = $1`

	row := r.Conn.QueryRow(query, id)

	var employee entity.Employee
	err := row.Scan(&employee.Id, &employee.Username, &employee.FirstName, &employee.LastName, &employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &employee, nil
}
