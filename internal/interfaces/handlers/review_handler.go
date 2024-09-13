package handlers

import (
	"errors"
	"net/http"
	"tenders/internal/application/interfaces"
	"tenders/internal/utils"
	"tenders/internal/utils/common"
	"tenders/internal/utils/consts"
)

type ReviewHandler struct {
	service interfaces.ReviewService
}

func NewReviewHandler(service interfaces.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	if common.CheckForExtraParams(r, []string{"username", "bidFeedback"}) {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectParams)
		return
	}

	bidId, err := common.GetUUIDFromRequestPath(r, "bidId")
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectBidId)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.NoUsernameParamPresent)
		return
	}

	bidFeedback := r.URL.Query().Get("bidFeedback")
	if bidFeedback == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectFeedback)
		return
	}

	bid, err := h.service.SubmitFeedback(bidId, username, bidFeedback)
	if err != nil {
		if errors.Is(err, utils.BidNotExistsError) {
			common.RespondWithError(w, http.StatusNotFound, consts.BidNotExists)
		} else if errors.Is(err, utils.UnauthorizedAccessError) {
			common.RespondWithError(w, http.StatusForbidden, consts.InsufficientPermissions)
		} else if errors.Is(err, utils.UserNotExistsError) {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		}
		return
	}
	common.RespondOKWithJson(w, bid)
}

func (h *ReviewHandler) GetReviewsList(w http.ResponseWriter, r *http.Request) {
	if common.CheckForExtraParams(r, []string{"authorUsername", "requesterUsername", "limit", "offset"}) {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectParams)
		return
	}

	tenderId, err := common.GetUUIDFromRequestPath(r, "tenderId")
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectTenderId)
		return
	}

	authorUsername := r.URL.Query().Get("authorUsername")
	if authorUsername == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.NoAuthorUsernameParamPresent)
		return
	}

	requesterUsername := r.URL.Query().Get("requesterUsername")
	if requesterUsername == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.NoRequesterUsernameParamPresent)
		return
	}

	limit, offset, err := common.GetPaginationParams(r)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectLimitOffsetParams)
		return
	}

	reviews, err := h.service.FindAllReviewsByBidAuthor(tenderId, authorUsername, requesterUsername, limit, offset)
	if err != nil {
		if errors.Is(err, utils.UserNotExistsError) {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else if errors.Is(err, utils.TenderNotExistsError) {
			common.RespondWithError(w, http.StatusNotFound, consts.TenderNotExists)
		} else if errors.Is(err, utils.BidForTenderNotExistsError) {
			common.RespondWithError(w, http.StatusBadRequest, consts.BidForTenderNotExistsError)
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		}
		return
	}

	if len(reviews) == 0 {
		common.RespondWithError(w, http.StatusNotFound, consts.ReviewsNotFound)
		return
	}

	common.RespondOKWithJson(w, reviews)
}
