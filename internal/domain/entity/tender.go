package entity

import (
	"github.com/google/uuid"
	"time"
)

type ServiceType string

const (
	CONSTRUCTION string = "Construction"
	DELIVERY     string = "Delivery"
	MANUFACTURE  string = "Manufacture"
)

type Tender struct {
	Id             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ServiceType    string    `json:"service_type"`
	Status         string    `json:"status"`
	Version        int       `json:"version"`
	OrganizationID uuid.UUID `json:"-"`
	CreatorID      uuid.UUID `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}
