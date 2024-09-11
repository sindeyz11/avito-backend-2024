package common

import (
	"encoding/json"
	"net/http"
	"tenders/internal/domain/dto"
	"tenders/internal/utils/consts"
)

func RespondWithError(w http.ResponseWriter, statusCode int, errorMsg string) {
	j, err := json.Marshal(dto.ErrorResponse{Reason: errorMsg})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	_, _ = w.Write(j)
}

func RespondWithJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		RespondWithError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
		return
	}
	w.WriteHeader(statusCode)
}
