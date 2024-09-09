package handlers

import (
	"encoding/json"
	"net/http"
	"tenders/internal/application/interfaces"
	"tenders/internal/domain/dto"
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

func (c *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, 405, consts.MethodNotAllowed)
		return
	}

	tenderRequest := dto.TenderRequest{}

	err := json.NewDecoder(r.Body).Decode(&tenderRequest)
	if err != nil {
		HandleError(w, 400, consts.IncorrectRequestBody)
		return
	}

	tender, err := c.service.CreateTender(&tenderRequest)

	if err != nil {
		var statusCode int
		if err.Error() == consts.UserNotExistsError {
			statusCode = 401
			HandleError(w, statusCode, consts.UserNotExists)
		} else {
			statusCode = 400
			HandleError(w, statusCode, err.Error())
		}
		return
	}

	err = json.NewEncoder(w).Encode(tender)
	w.WriteHeader(http.StatusOK)
}
