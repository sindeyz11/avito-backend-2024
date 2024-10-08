package persistence

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"strings"
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
	insertQuery := `
        INSERT INTO tender (
            name, description, service_type, status,
            organization_id, creator_id, created_at, tender_id, version
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING created_at
    `
	created := tender.CreatedAt.ConvertToTime()

	err := r.Conn.QueryRow(insertQuery,
		tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID,
		tender.CreatorID, created, tender.TenderId, tender.Version,
	).Scan(&tender.CreatedAt)

	if err != nil {
		return nil, err
	}
	return tender, nil
}

// FindAllAvailableByOrganizationId находит список тендеров организцаии к который принадлежит работник
func (r *TenderRepo) FindAllAvailableByOrganizationId(organizationId uuid.UUID, limit, offset int) ([]entity.Tender, error) {
	tenders := []entity.Tender{}
	queryStr := `
		SELECT DISTINCT ON (t.tender_id, t.name) t.tender_id, t.name, t.description, t.service_type, t.status, t.version,
		       t.organization_id, t.creator_id, t.created_at
		FROM tender t
		WHERE t.organization_id = $1 AND t.version = (
			SELECT MAX(version)
			FROM tender
			WHERE tender_id = t.tender_id
		)
		ORDER BY t.name ASC LIMIT $2 OFFSET $3
	`

	rows, err := r.Conn.Query(queryStr, organizationId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tender entity.Tender
		err = rows.Scan(
			&tender.TenderId, &tender.Name, &tender.Description,
			&tender.ServiceType, &tender.Status, &tender.Version,
			&tender.OrganizationID, &tender.CreatorID, &tender.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		tenders = append(tenders, tender)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tenders, nil
}

func (r *TenderRepo) FindAllPublished(serviceTypes []string, limit, offset int) ([]entity.Tender, error) {
	queryStr := `SELECT tender_id, name, description, service_type, status,
       	version, organization_id, creator_id, created_at 
		FROM tender t
		WHERE status = 'Published' AND t.version = (
			SELECT MAX(version)
			FROM tender
			WHERE tender_id = t.tender_id
		)`

	argIndex := 1
	var queryArgs []interface{}

	// Добавление фильтрации по типу услуг, если есть фильтры
	if len(serviceTypes) > 0 {
		placeholders := make([]string, len(serviceTypes))
		for i, serviceType := range serviceTypes {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			queryArgs = append(queryArgs, serviceType)
			argIndex++
		}
		queryStr += fmt.Sprintf(" AND service_type IN (%s)", strings.Join(placeholders, ","))
	}

	// Добавление сортировки и лимита с оффсетом
	queryStr += fmt.Sprintf(" ORDER BY name ASC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)

	queryArgs = append(queryArgs, limit, offset)
	rows, err := r.Conn.Query(queryStr, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenders := []entity.Tender{}
	for rows.Next() {
		var tender entity.Tender
		err = rows.Scan(
			&tender.TenderId, &tender.Name, &tender.Description,
			&tender.ServiceType, &tender.Status, &tender.Version,
			&tender.OrganizationID, &tender.CreatorID, &tender.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		tenders = append(tenders, tender)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tenders, nil
}

func (r *TenderRepo) FindByTenderId(tenderId uuid.UUID) (*entity.Tender, error) {
	var tender entity.Tender
	queryStr := `
		SELECT tender_id, name, description, service_type, status, version,
		       organization_id, creator_id, created_at
		FROM tender
		WHERE tender_id = $1 ORDER BY version DESC LIMIT 1
	`
	err := r.Conn.QueryRow(queryStr, tenderId).Scan(
		&tender.TenderId, &tender.Name, &tender.Description,
		&tender.ServiceType, &tender.Status, &tender.Version,
		&tender.OrganizationID, &tender.CreatorID, &tender.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &tender, nil
}

func (r *TenderRepo) FindByTenderIdAndVersion(tenderId uuid.UUID, version int) (*entity.Tender, error) {
	var tender entity.Tender
	queryStr := `
		SELECT tender_id, name, description, service_type, status, version,
		       organization_id, creator_id, created_at
		FROM tender
		WHERE tender_id = $1 AND version = $2
	`
	err := r.Conn.QueryRow(queryStr, tenderId, version).Scan(
		&tender.TenderId, &tender.Name, &tender.Description,
		&tender.ServiceType, &tender.Status, &tender.Version,
		&tender.OrganizationID, &tender.CreatorID, &tender.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &tender, nil
}

func (r *TenderRepo) FindLatestVersionByTenderId(tenderId uuid.UUID) (int, error) {
	var version int
	query := `
		SELECT MAX(version)
		FROM tender
		WHERE tender_id = $1
	`
	err := r.Conn.QueryRow(query, tenderId).Scan(&version)
	if err != nil {
		return -734, err
	}
	return version, nil
}

func (r *TenderRepo) CheckEmployeeAccessToTender(employeeId uuid.UUID, tenderId uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM tender t
			INNER JOIN organization o ON t.organization_id = o.id
			INNER JOIN organization_responsible orp ON o.id = orp.organization_id
			WHERE orp.user_id = $1 AND t.tender_id = $2
		);
	`

	var hasAccess bool
	err := r.Conn.QueryRow(query, employeeId, tenderId).Scan(&hasAccess)
	if err != nil {
		return false, err
	}

	return hasAccess, nil
}
