package service

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"tenders/internal/application/interfaces"
	"tenders/internal/domain/entity"
	"tenders/internal/domain/repository"
	"tenders/internal/interfaces/dto/request"
	"tenders/internal/utils"
	"tenders/internal/utils/common/custom_types"
	"tenders/internal/utils/consts"
	"time"
)

type TenderService struct {
	tenderRepo       repository.TenderRepository
	employeeRepo     repository.EmployeeRepository
	organizationRepo repository.OrganizationRepository
}

func NewTenderService(
	tenderRepo repository.TenderRepository,
	employeeRepo repository.EmployeeRepository,
	organizationRepo repository.OrganizationRepository,
) interfaces.TenderService {
	return &TenderService{
		tenderRepo:       tenderRepo,
		employeeRepo:     employeeRepo,
		organizationRepo: organizationRepo,
	}
}

func (s *TenderService) FindAllPublished(serviceTypes []string, limit, offset int) ([]entity.Tender, error) {
	return s.tenderRepo.FindAllPublished(serviceTypes, limit, offset)
}

func (s *TenderService) FindAllAvailableByEmployeeUsername(username string, limit, offset int) ([]entity.Tender, error) {
	employeeId, err := s.employeeRepo.FindEmployeeIdByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.UserNotExistsError
		}
		return nil, err
	}

	organization, err := s.organizationRepo.FindByEmployeeId(employeeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entity.Tender{}, nil
		}
		return nil, err
	}

	return s.tenderRepo.FindAllAvailableByOrganizationId(organization.Id, limit, offset)
}

func (s *TenderService) specifyEmployeeVerificationError(username string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		_, userNotFoundErr := s.employeeRepo.FindEmployeeIdByUsername(username)
		if userNotFoundErr != nil {
			return utils.UserNotExistsError
		}
		return utils.UnauthorizedAccessError
	}
	return err
}

func (s *TenderService) updateTenderWithVersionIncr(tender *entity.Tender) (*entity.Tender, error) {
	tender.Version = tender.Version + 1
	return s.tenderRepo.Create(tender)
}

func (s *TenderService) updateTenderFromOldVersion(tender *entity.Tender) (*entity.Tender, error) {
	latestVersion, err := s.tenderRepo.FindLatestVersionByTenderId(tender.TenderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.TenderNotExistsError
		}
		return nil, err
	}
	tender.Version = latestVersion + 1
	return s.tenderRepo.Create(tender)
}

func (s *TenderService) Create(tenderRequest *request.TenderRequest) (*entity.Tender, error) {
	tender, err := tenderRequest.MapToTender()
	if err != nil {
		return nil, err
	}

	employeeId, err := s.employeeRepo.FindEmployeeIdByUsernameIfResponsibleForOrg(
		tenderRequest.CreatorUsername,
		tenderRequest.OrganizationID,
	)
	if err != nil {
		err = s.specifyEmployeeVerificationError(tenderRequest.CreatorUsername, err)
		return nil, err
	}

	if tender.TenderId == uuid.Nil {
		tender.TenderId = uuid.New()
		tender.Version = 1
		tender.Status = consts.TenderCreated
		tender.CreatorID = employeeId
		tender.CreatedAt = custom_types.RFC3339Time(time.Now())
	}

	return s.tenderRepo.Create(tender)
}

func (s *TenderService) GetStatusByTenderId(tenderId uuid.UUID, username string) (string, error) {
	tender, err := s.tenderRepo.FindByTenderId(tenderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", utils.TenderNotExistsError
		}
		return "", err
	}

	if username == "" {
		if tender.Status != consts.TenderPublished {
			return "", utils.UnauthorizedAccessError
		}
		return tender.Status, nil
	}

	_, err = s.employeeRepo.FindEmployeeIdByUsernameIfResponsibleForOrg(username, tender.OrganizationID)
	if err != nil {
		err = s.specifyEmployeeVerificationError(username, err)
		return "", err
	}

	return tender.Status, nil
}

func (s *TenderService) VerifyUserResponsibleForOrg(username string, organizationId uuid.UUID) (uuid.UUID, error) {
	return s.employeeRepo.FindEmployeeIdByUsernameIfResponsibleForOrg(username, organizationId)
}

func (s *TenderService) GetTenderVersion(tenderId uuid.UUID, version int) (*entity.Tender, error) {
	return s.tenderRepo.FindByTenderIdAndVersion(tenderId, version)
}

func (s *TenderService) UpdateStatus(tenderId uuid.UUID, status, username string) (*entity.Tender, error) {
	tender, err := s.tenderRepo.FindByTenderId(tenderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.TenderNotExistsError
		}
		return nil, err
	}

	_, err = s.employeeRepo.FindEmployeeIdByUsernameIfResponsibleForOrg(username, tender.OrganizationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, userNotFoundErr := s.employeeRepo.FindEmployeeIdByUsername(username)
			if userNotFoundErr != nil {
				return nil, utils.UserNotExistsError
			}
			return nil, utils.UnauthorizedAccessError
		}
		return nil, err
	}

	tender.Status = status
	updatedTender, err := s.updateTenderWithVersionIncr(tender)
	if err != nil {
		return nil, err
	}

	return updatedTender, nil
}

func (s *TenderService) FindByTenderId(tenderId uuid.UUID) (*entity.Tender, error) {
	tender, err := s.tenderRepo.FindByTenderId(tenderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.TenderNotExistsError
		}
		return nil, err
	}
	return tender, nil
}

func (s *TenderService) EditTender(tenderId uuid.UUID, username string, updateRequest *request.EditTenderRequest) (*entity.Tender, error) {
	tender, err := s.tenderRepo.FindByTenderId(tenderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.TenderNotExistsError
		}
		return nil, err
	}

	_, err = s.employeeRepo.FindEmployeeIdByUsernameIfResponsibleForOrg(username, tender.OrganizationID)
	if err != nil {
		return nil, s.specifyEmployeeVerificationError(username, err)
	}

	if err = updateRequest.UpdateTender(tender); err != nil {
		return nil, err
	}

	return s.updateTenderWithVersionIncr(tender)
}

func (s *TenderService) RollbackTender(tenderId uuid.UUID, version int, username string) (*entity.Tender, error) {
	tender, err := s.GetTenderVersion(tenderId, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.TenderNotExistsError
		}
		return nil, err
	}

	_, err = s.VerifyUserResponsibleForOrg(username, tender.OrganizationID)
	if err != nil {
		return nil, s.specifyEmployeeVerificationError(username, err)
	}

	return s.updateTenderFromOldVersion(tender)
}
