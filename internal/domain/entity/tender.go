package entity

import (
	"github.com/google/uuid"
	"tenders/internal/utils/common/custom_types"
	"tenders/internal/utils/consts"
)

type Tender struct {
	Id             uuid.UUID                `json:"-"`
	TenderId       uuid.UUID                `json:"id"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	ServiceType    string                   `json:"service_type"`
	Status         string                   `json:"status"`
	Version        int                      `json:"version"`
	OrganizationID uuid.UUID                `json:"-"`
	CreatorID      uuid.UUID                `json:"-"`
	CreatedAt      custom_types.RFC3339Time `json:"created_at"`
} // По хорошему надо было добавить UpdatedAt, но т.к. он нигде не отдается - решил не добавлять

var ValidServiceTypes = map[string]bool{
	consts.Construction: true,
	consts.Delivery:     true,
	consts.Manufacture:  true,
}

var ValidTenderStatuses = map[string]bool{
	consts.TenderCreated:   true,
	consts.TenderPublished: true,
	consts.TenderClosed:    true,
}
