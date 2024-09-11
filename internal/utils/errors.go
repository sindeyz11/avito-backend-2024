package utils

import (
	"errors"
	"strings"
)

var ErrUnauthorizedAccess = errors.New("unauthorized access")
var ErrElementNotExist = errors.New("element not exist")

func NewValidationError(errorFields []string) error {
	return errors.New("Неправильно заполнены поля: " + strings.Join(errorFields, ", "))
}
