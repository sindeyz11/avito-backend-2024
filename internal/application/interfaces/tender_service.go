package interfaces

import (
	"github.com/google/uuid"
	"tenders/internal/domain/entity"
	"tenders/internal/interfaces/dto/request"
)

type TenderService interface {
	Create(tenderRequest *request.TenderRequest) (*entity.Tender, error)
	FindAllPublished(serviceTypes []string, limit, offset int) ([]entity.Tender, error)
	FindAllAvailableByEmployeeUsername(username string, limit, offset int) ([]entity.Tender, error)
	GetStatusByTenderId(id uuid.UUID, username string) (string, error)
	UpdateStatus(tenderId uuid.UUID, status, username string) (*entity.Tender, error)
	FindByTenderId(tenderId uuid.UUID) (*entity.Tender, error)
	VerifyUserResponsibleForOrg(username string, organizationId uuid.UUID) (uuid.UUID, error)
	GetTenderVersion(tenderId uuid.UUID, version int) (*entity.Tender, error)
	EditTender(tenderId uuid.UUID, username string, updateRequest *request.EditTenderRequest) (*entity.Tender, error)
	RollbackTender(tenderId uuid.UUID, version int, username string) (*entity.Tender, error)
}
