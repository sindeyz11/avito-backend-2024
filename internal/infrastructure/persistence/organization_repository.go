package persistence

import (
	"database/sql"
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/domain/repository"
)

type OrganizationRepo struct {
	Conn *sql.DB
}

func NewOrganizationRepository(conn *sql.DB) *OrganizationRepo {
	return &OrganizationRepo{Conn: conn}
}

var _ repository.OrganizationRepository = &OrganizationRepo{}

func (r *OrganizationRepo) FindById(id uuid.UUID) (*entity.Organization, error) {
	query := `
		SELECT id, name, description, type, created_at, updated_at
		FROM organization
		WHERE id = $1
	`

	var org entity.Organization
	err := r.Conn.QueryRow(query, id).Scan(
		&org.Id,
		&org.Name,
		&org.Description,
		&org.Type,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &org, nil
}

func (r *OrganizationRepo) FindByEmployeeId(employeeId uuid.UUID) (*entity.Organization, error) {
	query := `
		SELECT o.id, o.name, o.description, o.type, o.created_at, o.updated_at
		FROM organization o
		JOIN organization_responsible ore ON o.id = ore.organization_id
		WHERE ore.user_id = $1
	`

	var org entity.Organization
	err := r.Conn.QueryRow(query, employeeId).Scan(
		&org.Id,
		&org.Name,
		&org.Description,
		&org.Type,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &org, nil
}
