package service

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"tenders/internal/application/interfaces"
	"tenders/internal/domain/entity"
	"tenders/internal/domain/repository"
	"tenders/internal/utils"
	"tenders/internal/utils/common/custom_types"
	"time"
)

type ReviewService struct {
	bidRepo          repository.BidRepository
	employeeRepo     repository.EmployeeRepository
	organizationRepo repository.OrganizationRepository
	tenderRepo       repository.TenderRepository
	reviewRepo       repository.ReviewRepository
}

func NewReviewService(
	employeeRepo repository.EmployeeRepository,
	organizationRepo repository.OrganizationRepository,
	bidRepo repository.BidRepository,
	tenderRepo repository.TenderRepository,
	reviewRepo repository.ReviewRepository,
) interfaces.ReviewService {
	return &ReviewService{
		organizationRepo: organizationRepo,
		bidRepo:          bidRepo,
		tenderRepo:       tenderRepo,
		employeeRepo:     employeeRepo,
		reviewRepo:       reviewRepo,
	}
}

func (s *ReviewService) SubmitFeedback(bidId uuid.UUID, username string, feedback string) (*entity.Bid, error) {
	bid, err := s.bidRepo.FindByBidId(bidId)
	if err != nil {
		if errors.Is(err, utils.BidNotExistsError) {
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

	if hasAccess, _ := s.tenderRepo.CheckEmployeeAccessToTender(employeeId, bid.TenderId); !hasAccess {
		return nil, utils.UnauthorizedAccessError
	}

	review := &entity.Review{
		Id:          uuid.New(),
		BidId:       bidId,
		Description: feedback,
		CreatedAt:   custom_types.RFC3339Time(time.Now()),
	}

	if _, err = s.reviewRepo.Create(review); err != nil {
		return nil, err
	}
	return bid, nil
}

func (s *ReviewService) specifyEmployeeVerificationError(username string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		_, userNotFoundErr := s.employeeRepo.FindEmployeeIdByUsername(username)
		if userNotFoundErr != nil {
			return utils.UserNotExistsError
		}
		return utils.UnauthorizedAccessError
	}
	return err
}

func (s *ReviewService) FindAllReviewsByBidAuthor(
	tenderId uuid.UUID, authorUsername string, requesterUsername string, limit, offset int,
) ([]entity.Review, error) {
	authorId, err := s.employeeRepo.FindEmployeeIdByUsername(authorUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.UserNotExistsError
		}
		return nil, err
	}

	tender, err := s.tenderRepo.FindByTenderId(tenderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.TenderNotExistsError
		}
		return nil, err
	}

	_, err = s.bidRepo.FindByAuthorAndTender(authorId, tenderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.BidForTenderNotExistsError
		}
		return nil, err
	}

	_, err = s.employeeRepo.FindEmployeeIdByUsernameIfResponsibleForOrg(requesterUsername, tender.OrganizationID)
	if err != nil {
		err = s.specifyEmployeeVerificationError(requesterUsername, err)
		return nil, err
	}

	reviews, err := s.reviewRepo.FindAllByBidAuthor(authorId, limit, offset)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
