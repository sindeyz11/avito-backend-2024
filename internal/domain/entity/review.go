package entity

import (
	"github.com/google/uuid"
	"tenders/internal/utils/common/custom_types"
)

type Review struct {
	Id          uuid.UUID                `json:"id"`
	BidId       uuid.UUID                `json:"-"`
	Description string                   `json:"description"`
	CreatedAt   custom_types.RFC3339Time `json:"created_at"`
}
