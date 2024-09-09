package repository_impl

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/domain/repository"
)

type TenderRepo struct {
	Conn *sql.DB
}

func NewTenderRepository(conn *sql.DB) *TenderRepo {
	return &TenderRepo{Conn: conn}
}

var _ repository.TenderRepository = &TenderRepo{}

func (r *TenderRepo) Create(tender *entity.Tender) (*entity.Tender, error) {
	fmt.Println("tender repository hit")
	insertQuery := `
        INSERT INTO tender (name, description, service_type, status, organization_id, creator_id, created_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
    `
	err := r.Conn.QueryRow(insertQuery,
		tender.Name, tender.Description, tender.ServiceType,
		tender.Status, tender.OrganizationID, tender.CreatorID, tender.CreatedAt,
	).Scan(&tender.Id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return tender, nil
}

// FindByEmployeeUsername находит тендеры по имени пользователя, который их создал
func (r *TenderRepo) FindByEmployeeUsername(username string) (*entity.Tender, error) {
	var tender entity.Tender

	query := `
		SELECT t.id, t.name, t.description, t.service_type, t.status, t.organization_id, t.creator_id, t.created_at
		FROM tender t
		JOIN employee e ON t.creator_id = e.id
		WHERE e.username = $1
	`

	err := r.Conn.QueryRow(query, username).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationID,
		&tender.CreatorID,
		&tender.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no tenders found for this username")
		}
		return nil, err
	}

	return &tender, nil
}

func (r *TenderRepo) FindAll() ([]*entity.Tender, error) {
	query := `
		SELECT id, name, description, service_type, status, organization_id, creator_id, created_at
		FROM tender
	`
	rows, err := r.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []*entity.Tender
	for rows.Next() {
		var tender entity.Tender
		err = rows.Scan(
			&tender.Id,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.OrganizationID,
			&tender.CreatorID,
			&tender.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tenders, nil
}

func (r *TenderRepo) Update(tender *entity.Tender) (*entity.Tender, error) {
	query := `
		UPDATE tender 
		SET name = $1, description = $2, service_type = $3, status = $4, organization_id = $5
		WHERE id = $7
		RETURNING id, name, description, service_type, status, organization_id, creator_id, created_at
	`

	err := r.Conn.QueryRow(query,
		tender.Name,
		tender.Description,
		tender.ServiceType,
		tender.Status,
		tender.OrganizationID,
		tender.Id,
	).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationID,
		&tender.CreatorID,
		&tender.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return tender, nil
}

func (r *TenderRepo) Delete(id uuid.UUID) error {
	query := `DELETE FROM tender WHERE id = $1`

	result, err := r.Conn.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no tender found with the given id")
	}

	return nil
}
