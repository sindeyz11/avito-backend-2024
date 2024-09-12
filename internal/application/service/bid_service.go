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
	"time"
)

type BidService struct {
	employeeRepo repository.EmployeeRepository
	tenderRepo   repository.TenderRepository
	bidRepo      repository.BidRepository
}

func NewBidService(
	bidRepo repository.BidRepository, tenderRepo repository.TenderRepository, employeeRepo repository.EmployeeRepository,
) interfaces.BidService {
	return &BidService{
		bidRepo:      bidRepo,
		tenderRepo:   tenderRepo,
		employeeRepo: employeeRepo,
	}
}

func (s *BidService) Create(request *request.BidRequest) (*entity.Bid, error) {
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

	if request.AuthorType == entity.USER {
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

	if bid.BidId == uuid.Nil {
		bid.BidId = uuid.New()
		bid.Version = 1
	}
	bid.TenderVersion = tender.Version
	bid.CreatedAt = time.Now().Format(time.RFC3339)
	bid.Status = entity.CREATED

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

	return s.bidRepo.FindAllByEmployeeId(employeeId, limit, offset)
}

func (s *BidService) FindAllByTenderId(tenderId uuid.UUID, username string, limit, offset int) ([]*entity.Bid, error) {
	//TODO implement me
	panic("implement me")
}
