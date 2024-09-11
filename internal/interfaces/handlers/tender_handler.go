package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"tenders/internal/application/interfaces"
	"tenders/internal/domain/dto"
	"tenders/internal/utils"
	"tenders/internal/utils/common"
	"tenders/internal/utils/consts"
)

type TenderHandler struct {
	service interfaces.TenderService
}

func NewTenderHandler(
	service interfaces.TenderService,
) *TenderHandler {
	return &TenderHandler{
		service: service,
	}
}

func (h *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	var tenderRequest dto.TenderRequest
	err := json.NewDecoder(r.Body).Decode(&tenderRequest)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectRequestBody)
		return
	}

	tender, err := h.service.Create(&tenderRequest)
	if err != nil {
		// TODO в первом ифе может быть ситуация когда ввели рандомную оргу
		if err.Error() == consts.CannotFindUserError {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else if strings.HasPrefix(err.Error(), "Неправильно") {
			common.RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(tender); err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

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

	// TODO только published
	tenders, err := h.service.FindAll(serviceTypeFilter, limit, offset)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.UnknownBDError)
		return
	}

	common.RespondWithJson(w, http.StatusOK, tenders)
}

func (h *TenderHandler) GetAllTendersByUsername(w http.ResponseWriter, r *http.Request) {
	employeeUsername := r.URL.Query().Get("username")
	if employeeUsername == "" {
		common.RespondWithError(w, http.StatusBadRequest, consts.NoUsernameParamPresent)
		return
	}

	limit, offset, err := common.GetPaginationParams(r)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectLimitOffsetParams)
		return
	}

	tenders, err := h.service.FindAllByEmployeeUsername(employeeUsername, limit, offset)
	if err != nil {
		if errors.Is(err, utils.ErrElementNotExist) {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(tenders); err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
		return
	}
	w.WriteHeader(http.StatusOK)
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
		if errors.Is(err, sql.ErrNoRows) {
			common.RespondWithError(w, http.StatusNotFound, consts.TenderNotExists)
		} else if errors.Is(err, utils.ErrUnauthorizedAccess) {
			common.RespondWithError(w, http.StatusForbidden, consts.StatusForbidden)
		} else if err.Error() == consts.UserNotExistsError {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError)
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
		} else if errors.Is(err, utils.ErrUnauthorizedAccess) {
			common.RespondWithError(w, http.StatusForbidden, consts.StatusForbidden)
		} else if err.Error() == consts.UserNotExistsError {
			common.RespondWithError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(tender); err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
	}
	w.WriteHeader(http.StatusOK)
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

	var updateRequest dto.EditTenderRequest
	if err = json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, consts.IncorrectRequestBody)
		return
	}

	tender, err := h.service.FindByTenderId(tenderId)
	if err != nil {
		if errors.Is(err, utils.ErrElementNotExist) {
			common.RespondWithError(w, http.StatusNotFound, consts.TenderNotExists)
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.UnknownBDError)
		}
		return
	}

	_, err = h.service.VerifyUserResponsibleForOrg(username, tender.OrganizationID)
	if err != nil {
		if errors.Is(err, utils.ErrUnauthorizedAccess) {
			http.Error(w, "Forbidden: you are not responsible for this organization", http.StatusForbidden)
		} else if errors.Is(err, utils.ErrElementNotExist) {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			// todo fix
			common.RespondWithError(w, http.StatusInternalServerError, consts.UnknownBDError)
		}
		return
	}

	if err = updateRequest.MapToTender(tender); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedTender, err := h.service.UpdateTenderWithVersionIncr(tender)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.UnknownBDError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(updatedTender); err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
	}
	w.WriteHeader(http.StatusOK)
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

	tender, err := h.service.GetTenderByVersion(tenderId, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			common.RespondWithError(w, http.StatusNotFound, consts.TenderOrVersionNotExists)
		} else {
			common.RespondWithError(w, http.StatusInternalServerError, consts.UnknownBDError)
		}
		return
	}

	_, err = h.service.VerifyUserResponsibleForOrg(username, tender.OrganizationID)
	if err != nil {
		if errors.Is(err, utils.ErrUnauthorizedAccess) {
			http.Error(w, "Forbidden: you are not responsible for this organization", http.StatusForbidden)
		} else if errors.Is(err, utils.ErrElementNotExist) {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			// todo fix
			common.RespondWithError(w, http.StatusInternalServerError, consts.UnknownBDError)
		}
		return
	}

	updatedTender, err := h.service.UpdateTenderFromOldVersion(tender)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.UnknownBDError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(updatedTender); err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
	}
	w.WriteHeader(http.StatusOK)
}
