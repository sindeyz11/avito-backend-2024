package response

import (
	"strings"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

func NewValidationErrorResponse(errors []string) ErrorResponse {
	reason := "Неправильно заполнены поля: " + strings.Join(errors, ", ")
	return ErrorResponse{
		Reason: reason,
	}
}
