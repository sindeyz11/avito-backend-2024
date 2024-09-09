package persistence

import (
	"database/sql"
	"tenders/internal/domain/repository"
	"tenders/internal/infrastructure/persistence/repository_impl"
)

type Repositories struct {
	TenderRepo   repository.TenderRepository
	EmployeeRepo repository.EmployeeRepository
	Db           *sql.DB
}

func NewRepositories(conn *sql.DB) *Repositories {
	return &Repositories{
		TenderRepo:   repository_impl.NewTenderRepository(conn),
		EmployeeRepo: repository_impl.NewEmployeeRepository(conn),
		Db:           conn,
	}
}

func (r *Repositories) Close() error {
	return r.Db.Close()
}
