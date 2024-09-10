package utils

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"tenders/internal/domain/entity"
)

func NewValidationError(errorFields []string) error {
	return errors.New("Неправильно заполнены поля: " + strings.Join(errorFields, ", "))
}

func GetPaginationParams(r *http.Request) (int, int, error) {
	query := r.URL.Query()
	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")

	// Значения по умолчанию
	limit := 5
	offset := 0

	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			return 0, 0, errors.New("invalid limit parameter")
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return 0, 0, errors.New("invalid offset parameter")
		}
	}

	return limit, offset, nil
}

func GetServiceTypeFilter(r *http.Request) ([]string, error) {
	query := r.URL.Query()
	serviceTypeFilter := query["service_type"] // Массив фильтров по типу услуг

	for _, serviceType := range serviceTypeFilter {
		if !entity.ValidServiceTypes[serviceType] {
			return nil, errors.New("указан некорректный service_type: " + serviceType)
		}
	}

	return serviceTypeFilter, nil
}

func ValidateTenderStatus(status string) error {
	if !entity.ValidTenderStatuses[status] {
		return errors.New("указан некорректный статус тендера: " + status)
	}
	return nil
}
