package utils

import (
	"errors"
	"strings"
)

var ErrorUnauthorizedAccess = errors.New("unauthorized access")
var ErrorElementNotExist = errors.New("element not exist")

func NewValidationError(errorFields []string) error {
	return errors.New("Неправильно заполнены поля: " + strings.Join(errorFields, ", "))
}
