package entity

import (
	"github.com/google/uuid"
	"tenders/internal/utils/common/custom_types"
	"tenders/internal/utils/consts"
)

type Bid struct {
	Id            uuid.UUID                `json:"-"`
	BidId         uuid.UUID                `json:"id"`
	Name          string                   `json:"name"`
	Description   string                   `json:"-"`
	TenderId      uuid.UUID                `json:"-"`
	TenderVersion int                      `json:"-"`
	Status        string                   `json:"status"`
	AuthorType    string                   `json:"authorType"`
	AuthorId      uuid.UUID                `json:"authorId"`
	Version       int                      `json:"version"`
	CreatedAt     custom_types.RFC3339Time `json:"createdAt"`
}

var ValidBidStatuses = map[string]bool{
	consts.BidCreated:   true,
	consts.BidPublished: true,
	consts.BidCanceled:  true,
}
