package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"tenders/internal/application/interfaces"
	"tenders/internal/interfaces/dto/request"
	"tenders/internal/utils"
	"tenders/internal/utils/common"
	"tenders/internal/utils/consts"
)

type BidHandler struct {
	service interfaces.BidService
}

func NewBidHandler(service interfaces.BidService) *BidHandler {
	return &BidHandler{service: service}
}

func (h *BidHandler) CreateBid(w http.ResponseWriter, r *http.Request) {
	var bidRequest request.BidRequest
	if err := json.NewDecoder(r.Body).Decode(&bidRequest); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectRequestBody)
		return
	}

	bid, err := h.service.CreateNewBid(&bidRequest)
	if err != nil {
		if errors.Is(err, utils.ElementNotExistsError) {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserOrOrgNotExists)
		} else if errors.Is(err, utils.TenderNotExistsError) {
			common.RespondWithError(w, http.StatusNotFound, consts.TenderNotExists)
		} else if strings.HasPrefix(err.Error(), "Неправильно") {
			common.RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		}
		return
	}

	common.RespondOKWithJson(w, bid)
}

func (h *BidHandler) GetAllBidsByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.NoUsernameParamPresent)
		return
	}

	limit, offset, err := common.GetPaginationParams(r)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectLimitOffsetParams)
		return
	}

	tenders, err := h.service.FindAllByEmployeeUsername(username, limit, offset)
	if err != nil {
		if errors.Is(err, utils.UserNotExistsError) {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		}
		return
	}

	common.RespondOKWithJson(w, tenders)
}

func (h *BidHandler) GetAllBidsByTender(w http.ResponseWriter, r *http.Request) {
	tenderId, err := common.GetUUIDFromRequestPath(r, "tenderId")
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectTenderId)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.NoUsernameParamPresent)
		return
	}

	limit, offset, err := common.GetPaginationParams(r)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectLimitOffsetParams)
		return
	}

	bids, err := h.service.FindAllByTenderId(tenderId, username, limit, offset)
	if err != nil {
		if errors.Is(err, utils.TenderNotExistsError) {
			common.RespondWithError(w, http.StatusNotFound, consts.TenderNotExists)
		} else if errors.Is(err, utils.UnauthorizedAccessError) {
			common.RespondWithError(w, http.StatusForbidden, consts.InsufficientPermissions)
		} else if errors.Is(err, utils.UserNotExistsError) {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		}
		return
	}

	common.RespondOKWithJson(w, bids)
}

func (h *BidHandler) GetBidStatusById(w http.ResponseWriter, r *http.Request) {
	bidIdStr := r.PathValue("bidId")
	bidId, err := uuid.Parse(bidIdStr)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectTenderId)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.NoUsernameParamPresent)
		return
	}

	status, err := h.service.GetStatusByBidId(bidId, username)
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

	w.Header().Set("Content-Type", "text/plain")
	if _, err = w.Write([]byte(status)); err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
	}
}
func (h *BidHandler) UpdateBidStatusById(w http.ResponseWriter, r *http.Request) {
	bidId, err := common.GetUUIDFromRequestPath(r, "bidId")
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectBidId)
		return
	}

	status := r.URL.Query().Get("status")
	err = common.ValidateBidStatus(status)
	if status == "" || err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectStatus)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.NoUsernameParamPresent)
		return
	}

	bid, err := h.service.UpdateStatus(bidId, status, username)
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

func (h *BidHandler) EditBid(w http.ResponseWriter, r *http.Request) {
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

	var updateRequest request.EditBidRequest
	if err = json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectRequestBody)
		return
	}

	updatedBid, err := h.service.EditBid(bidId, username, &updateRequest)
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
	common.RespondOKWithJson(w, updatedBid)
}

func (h *BidHandler) RollbackBid(w http.ResponseWriter, r *http.Request) {
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

	version, err := common.GetVersionFromRequestPath(r)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectVersion)
		return
	}

	updatedBid, err := h.service.RollbackBid(bidId, version, username)
	if err != nil {
		if errors.Is(err, utils.BidNotExistsError) {
			common.RespondWithError(w, http.StatusNotFound, consts.BidNotExists)
		} else if errors.Is(err, utils.VersionNotExistsError) {
			common.RespondWithError(w, http.StatusNotFound, consts.VersionNotExists)
		} else if errors.Is(err, utils.UnauthorizedAccessError) {
			common.RespondWithError(w, http.StatusForbidden, consts.InsufficientPermissions)
		} else if errors.Is(err, utils.UserNotExistsError) {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		}
		return
	}
	common.RespondOKWithJson(w, updatedBid)
}

func (h *BidHandler) SubmitDecision(w http.ResponseWriter, r *http.Request) {
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

	decision := r.URL.Query().Get("decision")
	if decision != consts.BidApproved && decision != consts.BidRejected {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectDecision)
		return
	}

	updatedBid, err := h.service.SubmitDecision(bidId, username, decision)
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
	common.RespondOKWithJson(w, updatedBid)
}
