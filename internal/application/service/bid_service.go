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

type BidService struct {
	employeeRepo     repository.EmployeeRepository
	organizationRepo repository.OrganizationRepository
	tenderRepo       repository.TenderRepository
	bidRepo          repository.BidRepository
}

func NewBidService(
	employeeRepo repository.EmployeeRepository,
	organizationRepo repository.OrganizationRepository,
	bidRepo repository.BidRepository,
	tenderRepo repository.TenderRepository,
) interfaces.BidService {
	return &BidService{
		organizationRepo: organizationRepo,
		bidRepo:          bidRepo,
		tenderRepo:       tenderRepo,
		employeeRepo:     employeeRepo,
	}
}

func (s *BidService) CreateNewBid(request *request.BidRequest) (*entity.Bid, error) {
	bid, err := request.MapToBid()
	if err != nil {
		return nil, err
	}

	tender, err := s.tenderRepo.FindByTenderId(request.TenderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.TenderNotExistsError
		}
		return nil, err
	}

	if request.AuthorType == consts.AuthorTypeUser {
		_, err = s.employeeRepo.FindById(bid.AuthorId)
	} else {
		_, err = s.employeeRepo.FindOrgById(bid.AuthorId)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ElementNotExistsError
		}
		return nil, err
	}

	bid.BidId = uuid.New()
	bid.Version = 1
	bid.TenderVersion = tender.Version
	bid.CreatedAt = custom_types.RFC3339Time(time.Now())
	bid.Status = consts.BidCreated

	return s.bidRepo.Create(bid)
}

func (s *BidService) FindAllByEmployeeUsername(username string, limit, offset int) ([]entity.Bid, error) {
	employeeId, err := s.employeeRepo.FindEmployeeIdByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.UserNotExistsError
		}
		return nil, err
	}

	organization, err := s.organizationRepo.FindByEmployeeId(employeeId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	return s.bidRepo.FindAllByEmployeeIdAndOrgId(employeeId, organization.Id, limit, offset)
}

func (s *BidService) specifyEmployeeVerificationError(username string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		_, userNotFoundErr := s.employeeRepo.FindEmployeeIdByUsername(username)
		if userNotFoundErr != nil {
			return utils.UserNotExistsError
		}
		return utils.UnauthorizedAccessError
	}
	return err
}

func (s *BidService) FindAllByTenderId(tenderId uuid.UUID, username string, limit, offset int) ([]entity.Bid, error) {
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

	bids, err := s.bidRepo.FindAllByTenderId(tenderId, limit, offset)
	if err != nil {
		return nil, err
	}
	return bids, nil
}

func (s *BidService) GetStatusByBidId(bidId uuid.UUID, username string) (string, error) {
	bid, err := s.bidRepo.FindByBidId(bidId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", utils.BidNotExistsError
		}
		return "", err
	}

	employeeId, err := s.employeeRepo.FindEmployeeIdByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", utils.UserNotExistsError
		}
		return "", err
	}

	// Если создал юзер
	if bid.AuthorType == consts.AuthorTypeUser {
		if bid.AuthorId != employeeId {
			return "", utils.UnauthorizedAccessError
		}
		return bid.Status, nil
	} else {
		// Если создала организация
		org, err := s.organizationRepo.FindByEmployeeId(employeeId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", utils.UnauthorizedAccessError
			}
			return "", err
		}
		if bid.AuthorId != org.Id {
			return "", utils.UnauthorizedAccessError
		}

		return bid.Status, nil
	}
}

func (s *BidService) UpdateStatus(bidId uuid.UUID, status string, username string) (*entity.Bid, error) {
	// Открываем транзакцию
	tx, err := s.bidRepo.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	bid, err := s.bidRepo.FindByBidId(bidId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.BidNotExistsError
		}
		return nil, err
	}

	employeeId, err := s.employeeRepo.FindEmployeeIdByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.UserNotExistsError
		}
		return nil, err
	}

	if bid.AuthorType == consts.AuthorTypeUser {
		if bid.AuthorId != employeeId {
			return nil, utils.UnauthorizedAccessError
		}
	} else {
		org, err := s.organizationRepo.FindByEmployeeId(employeeId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, utils.UnauthorizedAccessError
			}
			return nil, err
		}
		if bid.AuthorId != org.Id {
			return nil, utils.UnauthorizedAccessError
		}
	}

	if err = s.bidRepo.SaveHistoricalVersionTx(tx, bid); err != nil {
		return nil, err
	}

	bid.Status = status
	bid.Version += 1

	err = s.bidRepo.UpdateStatusAndVersionTx(tx, bidId, status, bid.Version)
	if err != nil {
		return nil, err
	}

	return bid, nil
}

func (s *BidService) EditBid(bidId uuid.UUID, username string, updateRequest *request.EditBidRequest) (*entity.Bid, error) {
	// Открываем транзакцию
	tx, err := s.bidRepo.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	bid, err := s.bidRepo.FindByBidIdTx(tx, bidId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.BidNotExistsError
		}
		return nil, err
	}

	employeeId, err := s.employeeRepo.FindEmployeeIdByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.UserNotExistsError
		}
		return nil, err
	}

	if bid.AuthorType == consts.AuthorTypeUser {
		if bid.AuthorId != employeeId {
			return nil, utils.UnauthorizedAccessError
		}
	} else {
		org, err := s.organizationRepo.FindByEmployeeId(employeeId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, utils.UnauthorizedAccessError
			}
			return nil, err
		}
		if bid.AuthorId != org.Id {
			return nil, utils.UnauthorizedAccessError
		}
	}
	editedBid, err := updateRequest.MapToBid(*bid)
	if err != nil {
		return nil, err
	}

	if err = s.bidRepo.SaveHistoricalVersionTx(tx, bid); err != nil {
		return nil, err
	}

	editedBid.Version += 1
	if err = s.bidRepo.UpdateBidTx(tx, &editedBid); err != nil {
		return nil, err
	}

	return &editedBid, nil
}

func (s *BidService) RollbackBid(bidId uuid.UUID, version int, username string) (*entity.Bid, error) {
	// Открываем транзакцию
	tx, err := s.bidRepo.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	currentBid, err := s.bidRepo.FindByBidIdTx(tx, bidId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.BidNotExistsError
		}
		return nil, err
	}

	employeeId, err := s.employeeRepo.FindEmployeeIdByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.UserNotExistsError
		}
		return nil, err
	}

	if currentBid.AuthorType == consts.AuthorTypeUser {
		if currentBid.AuthorId != employeeId {
			return nil, utils.UnauthorizedAccessError
		}
	} else {
		org, err := s.organizationRepo.FindByEmployeeId(employeeId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, utils.UnauthorizedAccessError
			}
			return nil, err
		}
		if currentBid.AuthorId != org.Id {
			return nil, utils.UnauthorizedAccessError
		}
	}

	historicalBid, err := s.bidRepo.FindVersionInHistoryTx(tx, bidId, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.VersionNotExistsError
		}
		return nil, err
	}

	if err = s.bidRepo.SaveHistoricalVersionTx(tx, currentBid); err != nil {
		return nil, err
	}

	historicalBid.Version = currentBid.Version + 1

	if err = s.bidRepo.UpdateBidTx(tx, historicalBid); err != nil {
		return nil, err
	}

	return historicalBid, nil
}

func (s *BidService) SubmitDecision(bidId uuid.UUID, username, decision string) (*entity.Bid, error) {
	bid, err := s.bidRepo.FindByBidId(bidId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.BidNotExistsError
		}
		return nil, err
	}

	tender, err := s.tenderRepo.FindByTenderId(bid.TenderId)
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

	if decision == consts.BidApproved {
		tender.Status = consts.TenderClosed
		tender.Version = tender.Version + 1
		_, err = s.tenderRepo.Create(tender)
		if err != nil {
			return nil, err
		}
	}

	return bid, nil
}
