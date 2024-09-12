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

	bid, err := h.service.Create(&bidRequest)
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

func (h *BidHandler) GetBidsForTender(w http.ResponseWriter, r *http.Request) {
	tenderIdStr := r.URL.Query().Get("tenderId")
	tenderId, err := uuid.Parse(tenderIdStr)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid tender ID")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Username parameter is required")
		return
	}

	limit, offset, err := common.GetPaginationParams(r)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid pagination parameters")
		return
	}

	bids, err := h.service.FindAllByTenderId(tenderId, username, limit, offset)
	if err != nil {
		return
	}

	common.RespondOKWithJson(w, bids)
}
