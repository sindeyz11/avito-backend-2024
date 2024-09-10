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
		HandleError(w, http.StatusBadRequest, consts.IncorrectRequestBody)
		return
	}

	tender, err := h.service.Create(&tenderRequest)
	if err != nil {
		// TODO в первом ифе может быть ситуация когда ввели рандомную оргу
		if err.Error() == consts.CannotFindUserError {
			HandleError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else if strings.HasPrefix(err.Error(), "Неправильно") {
			HandleError(w, http.StatusBadRequest, err.Error())
		} else {
			HandleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(tender); err != nil {
		HandleError(w, http.StatusInternalServerError, consts.InternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TenderHandler) GetAllTenders(w http.ResponseWriter, r *http.Request) {
	serviceTypeFilter, err := utils.GetServiceTypeFilter(r)
	if err != nil {
		HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	limit, offset, err := utils.GetPaginationParams(r)
	if err != nil {
		HandleError(w, http.StatusBadRequest, consts.IncorrectLimitOffsetParams)
		return
	}

	// TODO только published
	tenders, err := h.service.FindAll(serviceTypeFilter, limit, offset)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, consts.InternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(tenders); err != nil {
		HandleError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TenderHandler) GetAllTendersByUsername(w http.ResponseWriter, r *http.Request) {
	employeeUsername := r.URL.Query().Get("username")
	if employeeUsername == "" {
		HandleError(w, http.StatusBadRequest, consts.NoUsernameParamPresent)
		return
	}

	limit, offset, err := utils.GetPaginationParams(r)
	if err != nil {
		HandleError(w, http.StatusBadRequest, consts.IncorrectLimitOffsetParams)
		return
	}

	tenders, err := h.service.FindAllByEmployeeUsername(employeeUsername, limit, offset)
	if err != nil {
		if err.Error() == consts.UserNotExistsError {
			HandleError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else {
			HandleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(tenders); err != nil {
		HandleError(w, http.StatusInternalServerError, consts.InternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TenderHandler) GetTenderStatusById(w http.ResponseWriter, r *http.Request) {
	tenderIdStr := r.PathValue("tenderId")
	tenderId, err := uuid.Parse(tenderIdStr)
	if err != nil {
		HandleError(w, http.StatusBadRequest, consts.IncorrectTenderId)
		return
	}

	username := r.URL.Query().Get("username")

	status, err := h.service.GetStatusByTenderId(tenderId, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			HandleError(w, http.StatusNotFound, consts.TenderNotExists)
		} else if errors.Is(err, utils.ErrUnauthorizedAccess) {
			HandleError(w, http.StatusForbidden, consts.StatusForbidden)
		} else if err.Error() == consts.UserNotExistsError {
			HandleError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else {
			HandleError(w, http.StatusInternalServerError, consts.InternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(status)); err != nil {
		HandleError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
	}
}

func (h *TenderHandler) UpdateTenderStatusById(w http.ResponseWriter, r *http.Request) {
	tenderIdStr := r.PathValue("tenderId")
	tenderId, err := uuid.Parse(tenderIdStr)
	if err != nil {
		HandleError(w, http.StatusBadRequest, consts.IncorrectTenderId)
		return
	}

	status := r.URL.Query().Get("status")
	err = utils.ValidateTenderStatus(status)
	if status == "" || err != nil {
		HandleError(w, http.StatusBadRequest, consts.IncorrectTenderStatus)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		HandleError(w, http.StatusBadRequest, consts.NoUsernameParamPresent)
		return
	}

	tender, err := h.service.UpdateStatus(tenderId, status, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			HandleError(w, http.StatusNotFound, consts.TenderNotExists)
		} else if errors.Is(err, utils.ErrUnauthorizedAccess) {
			HandleError(w, http.StatusForbidden, consts.StatusForbidden)
		} else if err.Error() == consts.UserNotExistsError {
			HandleError(w, http.StatusUnauthorized, consts.UserNotExists)
		} else {
			HandleError(w, http.StatusInternalServerError, consts.InternalServerError+" "+err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(tender); err != nil {
		HandleError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
	}
}
