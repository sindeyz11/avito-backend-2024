package handlers

import (
	"database/sql"
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

type TenderHandler struct {
	service interfaces.TenderService
}

func NewTenderHandler(service interfaces.TenderService) *TenderHandler {
	return &TenderHandler{service: service}
}

func (h *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	var tenderRequest request.TenderRequest
	if err := common.DecodeAndValidateJSON(r.Body, &tenderRequest); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectRequestBody)
		return
	}

	if common.CheckForExtraParams(r, []string{}) {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectParams)
		return
	}

	tender, err := h.service.Create(&tenderRequest)
	if err != nil {
		if errors.Is(err, utils.UserNotExistsError) {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else if errors.Is(err, utils.UnauthorizedAccessError) {
			common.RespondWithError(w, http.StatusForbidden, consts.InsufficientPermissions)
		} else if strings.HasPrefix(err.Error(), "Неправильно") {
			common.RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		}
		return
	}

	common.RespondOKWithJson(w, tender)
}

func (h *TenderHandler) GetAllTenders(w http.ResponseWriter, r *http.Request) {
	if common.CheckForExtraParams(r, []string{"service_type", "limit", "offset"}) {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectParams)
		return
	}

	serviceTypeFilter, err := common.GetServiceTypeFilter(r)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	limit, offset, err := common.GetPaginationParams(r)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectLimitOffsetParams)
		return
	}

	tenders, err := h.service.FindAllPublished(serviceTypeFilter, limit, offset)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		return
	}

	common.RespondOKWithJson(w, tenders)
}

func (h *TenderHandler) GetAllTendersByUsername(w http.ResponseWriter, r *http.Request) {
	if common.CheckForExtraParams(r, []string{"username", "limit", "offset"}) {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectParams)
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

	tenders, err := h.service.FindAllAvailableByEmployeeUsername(username, limit, offset)
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

func (h *TenderHandler) GetTenderStatusById(w http.ResponseWriter, r *http.Request) {
	tenderIdStr := r.PathValue("tenderId")
	tenderId, err := uuid.Parse(tenderIdStr)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectTenderId)
		return
	}

	if common.CheckForExtraParams(r, []string{"username"}) {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectParams)
		return
	}
	username := r.URL.Query().Get("username")

	status, err := h.service.GetStatusByTenderId(tenderId, username)
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

	w.Header().Set("Content-Type", "text/plain")
	if _, err = w.Write([]byte(status)); err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TenderHandler) UpdateTenderStatusById(w http.ResponseWriter, r *http.Request) {
	tenderId, err := common.GetUUIDFromRequestPath(r, "tenderId")
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectTenderId)
		return
	}

	if common.CheckForExtraParams(r, []string{"username", "status"}) {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectParams)
		return
	}

	status := r.URL.Query().Get("status")
	err = common.ValidateTenderStatus(status)
	if status == "" || err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectStatus)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.NoUsernameParamPresent)
		return
	}

	tender, err := h.service.UpdateStatus(tenderId, status, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

	common.RespondOKWithJson(w, tender)
}

func (h *TenderHandler) EditTender(w http.ResponseWriter, r *http.Request) {
	tenderId, err := common.GetUUIDFromRequestPath(r, "tenderId")
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectTenderId)
		return
	}

	if common.CheckForExtraParams(r, []string{"username"}) {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectParams)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.NoUsernameParamPresent)
		return
	}

	var updateRequest request.EditTenderRequest
	if err := common.DecodeAndValidateJSON(r.Body, &updateRequest); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectRequestBody)
		return
	}

	updatedTender, err := h.service.EditTender(tenderId, username, &updateRequest)
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

	common.RespondOKWithJson(w, updatedTender)
}

func (h *TenderHandler) RollbackTender(w http.ResponseWriter, r *http.Request) {
	tenderId, err := common.GetUUIDFromRequestPath(r, "tenderId")
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectTenderId)
		return
	}

	if common.CheckForExtraParams(r, []string{"username", "version"}) {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectParams)
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

	updatedTender, err := h.service.RollbackTender(tenderId, version, username)
	if err != nil {
		if errors.Is(err, utils.TenderNotExistsError) {
			common.RespondWithError(w, http.StatusNotFound, consts.TenderOrVersionNotExists)
		} else if errors.Is(err, utils.UnauthorizedAccessError) {
			common.RespondWithError(w, http.StatusForbidden, consts.InsufficientPermissions)
		} else if errors.Is(err, utils.UserNotExistsError) {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		}
		return
	}

	common.RespondOKWithJson(w, updatedTender)
}
