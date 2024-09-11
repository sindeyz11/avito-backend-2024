package service

import (
	"errors"
	"github.com/google/uuid"
	"tenders/internal/application/interfaces"
	"tenders/internal/domain/dto"
	"tenders/internal/domain/entity"
	"tenders/internal/domain/repository"
	"tenders/internal/utils"
	"tenders/internal/utils/consts"
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

func (s *TenderService) Create(tenderRequest *dto.TenderRequest) (*entity.Tender, error) {
	tender, err := tenderRequest.MapToTender()
	if err != nil {
		return nil, err
	}

	employeeId, err := s.employeeRepo.FindEmployeeIdByUsernameIfResponsibleForOrg(
		tenderRequest.CreatorUsername,
		tenderRequest.OrganizationID,
	)
	if err != nil {
		return nil, err
	}

	if tender.TenderId == uuid.Nil {
		tender.TenderId = uuid.New()
		tender.Version = 1
	}
	tender.CreatorID = employeeId
	tender.CreatedAt = time.Now()

	return s.tenderRepo.Create(tender)
}

func (s *TenderService) FindAll(serviceTypes []string, limit, offset int) ([]entity.Tender, error) {
	return s.tenderRepo.FindAll(serviceTypes, limit, offset)
}

func (s *TenderService) FindAllByEmployeeUsername(username string, limit, offset int) ([]entity.Tender, error) {
	employeeId, err := s.employeeRepo.FindEmployeeIdByUsername(username)
	if err != nil {
		return nil, err
	}

	return s.tenderRepo.FindAllByEmployeeId(employeeId, limit, offset)
}

func (s *TenderService) GetStatusByTenderId(tenderId uuid.UUID, username string) (string, error) {
	tender, err := s.tenderRepo.FindByTenderId(tenderId)
	if err != nil {
		return "", err
	}

	if username == "" {
		if tender.Status != entity.PUBLISHED {
			return "", utils.ErrUnauthorizedAccess
		}
		return tender.Status, nil
	}

	_, err = s.employeeRepo.FindEmployeeIdByUsernameIfResponsibleForOrg(username, tender.OrganizationID)
	if err != nil {
		_, errEmployeeNotExists := s.employeeRepo.FindEmployeeIdByUsername(username)
		if errEmployeeNotExists != nil {
			return "", errEmployeeNotExists
		}
		return "", utils.ErrUnauthorizedAccess
	}

	return tender.Status, nil
}

func (s *TenderService) UpdateTenderWithVersionIncr(tender *entity.Tender) (*entity.Tender, error) {
	tender.Version = tender.Version + 1
	return s.tenderRepo.Create(tender)
}

func (s *TenderService) UpdateTenderFromOldVersion(tender *entity.Tender) (*entity.Tender, error) {
	latestVersion, err := s.tenderRepo.FindLatestVersionByTenderId(tender.TenderId)
	if err != nil {
		return nil, err
	}
	tender.Version = latestVersion + 1
	return s.tenderRepo.Create(tender)
}

func (s *TenderService) UpdateStatus(tenderId uuid.UUID, status, username string) (*entity.Tender, error) {
	tender, err := s.tenderRepo.FindByTenderId(tenderId)
	if err != nil {
		return nil, err
	}

	_, err = s.employeeRepo.FindEmployeeIdByUsernameIfResponsibleForOrg(username, tender.OrganizationID)
	if err != nil {
		_, errEmployeeNotExists := s.employeeRepo.FindEmployeeIdByUsername(username)
		if errEmployeeNotExists != nil {
			return nil, errors.New(consts.UserNotExistsError)
		}
		return nil, utils.ErrUnauthorizedAccess
	}

	tender.Status = status
	updatedTender, err := s.UpdateTenderWithVersionIncr(tender)
	if err != nil {
		return nil, err
	}

	return updatedTender, nil
}

func (s *TenderService) FindByTenderId(tenderId uuid.UUID) (*entity.Tender, error) {
	tender, err := s.tenderRepo.FindByTenderId(tenderId)
	if err != nil {
		return nil, err
	}
	return tender, nil
}

func (s *TenderService) VerifyUserResponsibleForOrg(username string, organizationId uuid.UUID) (uuid.UUID, error) {
	return s.employeeRepo.FindEmployeeIdByUsernameIfResponsibleForOrg(username, organizationId)
}

func (s *TenderService) GetTenderByVersion(tenderId uuid.UUID, version int) (*entity.Tender, error) {
	return s.tenderRepo.FindByTenderIdAndVersion(tenderId, version)
}

// todo после рефакторинга
func (s *TenderService) RollbackTender(tenderId uuid.UUID, version int, username string) (*entity.Tender, error) {
	// Получаем тендер по ID и версии
	tender, err := s.GetTenderByVersion(tenderId, version)
	if err != nil {
		return nil, err
	}

	// Проверяем, что пользователь ответственен за организацию
	_, err = s.VerifyUserResponsibleForOrg(username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	// Обновляем тендер на основе старой версии
	return s.UpdateTenderFromOldVersion(tender)
}
