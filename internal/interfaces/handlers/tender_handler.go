package handlers

import (
	"database/sql"
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

type TenderHandler struct {
	service interfaces.TenderService
}

func NewTenderHandler(service interfaces.TenderService) *TenderHandler {
	return &TenderHandler{service: service}
}

func (h *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	var tenderRequest request.TenderRequest
	if err := json.NewDecoder(r.Body).Decode(&tenderRequest); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectRequestBody)
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

// TODO мб надо добавить фильтр по PUBLISHED
func (h *TenderHandler) GetAllTenders(w http.ResponseWriter, r *http.Request) {
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

	tenders, err := h.service.FindAll(serviceTypeFilter, limit, offset)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		return
	}

	common.RespondOKWithJson(w, tenders)
}

// TODO мб надо добавить фильтр по доступных организации пользователя
func (h *TenderHandler) GetAllTendersByUsername(w http.ResponseWriter, r *http.Request) {
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
	tenderId, err := common.GetTenderUUIDFromRequestPath(r)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectTenderId)
		return
	}

	status := r.URL.Query().Get("status")
	err = common.ValidateTenderStatus(status)
	if status == "" || err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectTenderStatus)
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
	tenderId, err := common.GetTenderUUIDFromRequestPath(r)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectTenderId)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.NoUsernameParamPresent)
		return
	}

	var updateRequest request.EditTenderRequest
	if err = json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
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
	tenderId, err := common.GetTenderUUIDFromRequestPath(r)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectTenderId)
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
