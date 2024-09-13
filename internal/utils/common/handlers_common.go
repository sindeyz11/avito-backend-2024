package common

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"tenders/internal/interfaces/dto/response"
	"tenders/internal/utils"
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

func DecodeAndValidateJSON(body io.Reader, v interface{}) error {
	var rawRequest map[string]interface{}
	if err := json.NewDecoder(body).Decode(&rawRequest); err != nil {
		return utils.IncorrectRequestBody
	}

	rawRequestBytes, err := json.Marshal(rawRequest)
	if err != nil {
		return utils.IncorrectRequestBody
	}

	dec := json.NewDecoder(strings.NewReader(string(rawRequestBytes)))
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok || strings.Contains(err.Error(), "unknown field") {
			return utils.IncorrectRequestBody
		}
		return utils.IncorrectRequestBody
	}

	return nil
}
