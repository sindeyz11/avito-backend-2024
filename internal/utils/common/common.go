package common

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"tenders/internal/domain/entity"
)

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

func ValidateBidStatus(status string) error {
	if !entity.ValidBidStatuses[status] {
		return errors.New("указан некорректный статус предложения: " + status)
	}
	return nil
}

func GetUUIDFromRequestPath(r *http.Request, pathValue string) (uuid.UUID, error) {
	uuidStr := r.PathValue(pathValue)
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func GetVersionFromRequestPath(r *http.Request) (int, error) {
	versionStr := r.PathValue("version")

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return -734, errors.New("невалидная версия")
	}

	return version, nil
}
