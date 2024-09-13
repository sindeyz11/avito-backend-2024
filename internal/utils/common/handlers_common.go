package common

import (
	"encoding/json"
	"net/http"
	"tenders/internal/interfaces/dto/response"
	"tenders/internal/utils/consts"
)

func RespondWithError(w http.ResponseWriter, statusCode int, errorMsg string) {
	j, err := json.Marshal(response.ErrorResponse{Reason: errorMsg})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	_, _ = w.Write(j)
}

func RespondOKWithJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		RespondWithError(w, http.StatusInternalServerError, consts.FailedToWriteResponse)
		return
	}
}

// CheckForExtraParams возвращает true, если есть лишние параметры
func CheckForExtraParams(r *http.Request, expectedParams []string) bool {
	// Создаем мапу для быстрого поиска ожидаемых параметров
	expectedParamsMap := make(map[string]struct{}, len(expectedParams))
	for _, param := range expectedParams {
		expectedParamsMap[param] = struct{}{}
	}

	// Извлекаем параметры запроса
	query := r.URL.Query()

	// Перебираем все параметры из запроса
	for param := range query {
		// Если параметр не входит в список ожидаемых, возвращаем true
		if _, found := expectedParamsMap[param]; !found {
			return true
		}
	}

	// Нет лишних параметров
	return false
}
