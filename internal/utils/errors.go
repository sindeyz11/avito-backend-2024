package utils

import (
	"errors"
	"strings"
)

var (
	UnauthorizedAccessError = errors.New("unauthorized access")
	DBError                 = errors.New("db error")

	ElementNotExistsError = errors.New("element not exist")
	TenderNotExistsError  = errors.New("tender not exist")
	BidNotExistsError     = errors.New("bid not exist")
	UserNotExistsError    = errors.New("user not exist")
)

func NewValidationError(errorFields []string) error {
	return errors.New("Неправильно заполнены поля: " + strings.Join(errorFields, ", "))
}
