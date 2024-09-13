package persistence

import (
	"database/sql"
	"tenders/internal/domain/repository"
)

type Repositories struct {
	TenderRepo       repository.TenderRepository
	EmployeeRepo     repository.EmployeeRepository
	OrganizationRepo repository.OrganizationRepository
	BidRepo          repository.BidRepository
	Db               *sql.DB
}

func NewRepositories(conn *sql.DB) *Repositories {
	return &Repositories{
		TenderRepo:       NewTenderRepository(conn),
		EmployeeRepo:     NewEmployeeRepository(conn),
		OrganizationRepo: NewOrganizationRepository(conn),
		BidRepo:          NewBidRepository(conn),
		Db:               conn,
	}
}

func (r *Repositories) Close() error {
	return r.Db.Close()
}
