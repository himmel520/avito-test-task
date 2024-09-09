package models

import "time"

// TenderStatus представляет статус тендера
type TenderStatus string

const (
	TenderStatusCreated   TenderStatus = "Created"
	TenderStatusPublished TenderStatus = "Published"
	TenderStatusClosed    TenderStatus = "Closed"
)

// TenderServiceType представляет вид услуги, к которой относится тендер
type TenderServiceType string

const (
	TenderServiceTypeConstruction TenderServiceType = "Construction"
	TenderServiceTypeDelivery     TenderServiceType = "Delivery"
	TenderServiceTypeManufacture  TenderServiceType = "Manufacture"
)

// TenderResponse представляет информацию о тендере
type TenderResponse struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	ServiceType    TenderServiceType `json:"serviceType"`
	Status         TenderStatus      `json:"status"`
	OrganizationID string            `json:"organizationId"`
	Version        int               `json:"version"`
	CreatedAt      time.Time         `json:"createdAt"`
}

// TenderCreate представляет данные для создания нового тендера
type TenderCreate struct {
	Name            string            `json:"name" validate:"required,max=100"`
	Description     string            `json:"description" validate:"required,max=500"`
	ServiceType     TenderServiceType `json:"serviceType" validate:"required,oneof=Construction Delivery Manufacture"`
	Status          TenderStatus      `json:"status" validate:"required,oneof=Created Published Closed"`
	OrganizationID  OrganizationID    `json:"organizationId" validate:"required,max=100"`
	CreatorUsername string            `json:"creatorUsername" validate:"required,max=50"`
}

// TenderEdit представляет данные для редактирования существующего тендера
type TenderEdit struct {
	Name        *string            `json:"name" validate:"omitempty,max=100"`
	Description *string            `json:"description" validate:"omitempty,max=500"`
	ServiceType *TenderServiceType `json:"serviceType" validate:"omitempty,oneof=Construction Delivery Manufacture"`
}

// OrganizationID представляет уникальный идентификатор организации
type OrganizationID string
