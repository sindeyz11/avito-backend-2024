package utils

import (
	"errors"
	"strings"
)

var (
	UnauthorizedAccessError = errors.New("unauthorized access")
	IncorrectRequestBody    = errors.New("incorrect request body")

	ElementNotExistsError      = errors.New("element not exist")
	TenderNotExistsError       = errors.New("tender not exist")
	BidNotExistsError          = errors.New("bid not exist")
	BidForTenderNotExistsError = errors.New("bid for tender not exist")
	UserNotExistsError         = errors.New("user not exist")
	VersionNotExistsError      = errors.New("version not exist")
)

func NewValidationError(errorFields []string) error {
	return errors.New("Неправильно заполнены поля: " + strings.Join(errorFields, ", "))
}
