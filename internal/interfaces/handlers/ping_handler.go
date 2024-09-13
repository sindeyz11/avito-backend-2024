package handlers

import (
	"net/http"
	"tenders/internal/utils/common"
	"tenders/internal/utils/consts"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte("ok")); err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
	}
}
