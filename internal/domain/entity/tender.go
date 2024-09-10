package entity

import (
	"github.com/google/uuid"
	"time"
)

type Tender struct {
	Id             uuid.UUID `json:"-"`
	TenderId       uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ServiceType    string    `json:"service_type"`
	Status         string    `json:"status"`
	Version        int       `json:"version"`
	OrganizationID uuid.UUID `json:"-"`
	CreatorID      uuid.UUID `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
} // По хорошему надо было добавить UpdatedAt, но т.к. он нигде не отдается - решил не добавлять

const (
	CONSTRUCTION string = "Construction"
	DELIVERY     string = "Delivery"
	MANUFACTURE  string = "Manufacture"
)

var ValidServiceTypes = map[string]bool{
	CONSTRUCTION: true,
	DELIVERY:     true,
	MANUFACTURE:  true,
}

const (
	CREATED   string = "Created"
	PUBLISHED string = "Published"
	CLOSED    string = "Closed"
)

var ValidTenderStatuses = map[string]bool{
	CREATED:   true,
	PUBLISHED: true,
	CLOSED:    true,
}
