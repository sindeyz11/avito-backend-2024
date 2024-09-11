package interfaces

import (
	"github.com/google/uuid"
	"tenders/internal/domain/dto"
	"tenders/internal/domain/entity"
)

type TenderService interface {
	Create(tenderRequest *dto.TenderRequest) (*entity.Tender, error)
	FindAll(serviceTypes []string, limit, offset int) ([]entity.Tender, error)
	FindAllByEmployeeUsername(username string, limit, offset int) ([]entity.Tender, error)
	GetStatusByTenderId(id uuid.UUID, username string) (string, error)
	UpdateTenderWithVersionIncr(tender *entity.Tender) (*entity.Tender, error)
	UpdateTenderFromOldVersion(tender *entity.Tender) (*entity.Tender, error)
	UpdateStatus(tenderId uuid.UUID, status, username string) (*entity.Tender, error)
	FindByTenderId(tenderId uuid.UUID) (*entity.Tender, error)
	VerifyUserResponsibleForOrg(username string, organizationId uuid.UUID) (uuid.UUID, error)
	GetTenderByVersion(tenderId uuid.UUID, version int) (*entity.Tender, error)
	RollbackTender(tenderId uuid.UUID, version int, username string) (*entity.Tender, error)
}
