package service

import (
	"tenders/internal/application/interfaces"
	"tenders/internal/domain/dto"
	"tenders/internal/domain/entity"
	"tenders/internal/domain/repository"
	"time"
)

type TenderService struct {
	tenderRepo   repository.TenderRepository
	employeeRepo repository.EmployeeRepository
}

func NewTenderService(
	tenderRepo repository.TenderRepository, employeeRepo repository.EmployeeRepository,
) interfaces.TenderService {
	return &TenderService{
		tenderRepo:   tenderRepo,
		employeeRepo: employeeRepo,
	}
}

func (s TenderService) CreateTender(tenderRequest *dto.TenderRequest) (*entity.Tender, error) {
	tender, err := tenderRequest.MapToTender()
	if err != nil {
		return nil, err
	}

	employeeId, err := s.employeeRepo.GetEmployeeIdByUsernameIfResponsibleForOrg(
		tenderRequest.CreatorUsername,
		tenderRequest.OrganizationID,
	)
	if err != nil {
		return nil, err
	}

	tender.CreatorID = employeeId
	tender.CreatedAt = time.Now()
	tender.Version = 1

	return s.tenderRepo.Create(tender)
}
