package entity

import (
	"github.com/google/uuid"
)

type Bid struct {
	Id            uuid.UUID `json:"-"`
	BidId         uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"-"`
	TenderId      uuid.UUID `json:"-"`
	TenderVersion int       `json:"-"`
	Status        string    `json:"status"`
	AuthorType    string    `json:"authorType"`
	AuthorId      uuid.UUID `json:"authorId"`
	Version       int       `json:"version"`
	CreatedAt     string    `json:"createdAt"`
}

const (
	ORGANIZATION string = "Organization"
	USER         string = "User"

	CANCELED string = "Canceled"
)
