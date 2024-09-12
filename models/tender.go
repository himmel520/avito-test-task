package models

import (
	"time"
)

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

// OrganizationID представляет id организации
type OrganizationID string

// TenderResponse представляет информацию о тендере
type TenderResponse struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	ServiceType    TenderServiceType `json:"serviceType"`
	Status         TenderStatus      `json:"status"`
	OrganizationID OrganizationID    `json:"organizationId"`
	Version        int               `json:"version"`
	CreatedAt      time.Time         `json:"createdAt"`
}

// TenderCreate представляет данные для создания нового тендера
type TenderCreate struct {
	Name            string            `json:"name" binding:"required,max=100"`
	Description     string            `json:"description" binding:"required,max=500"`
	ServiceType     TenderServiceType `json:"serviceType" binding:"required,oneof=Construction Delivery Manufacture"`
	OrganizationID  OrganizationID    `json:"organizationId" binding:"required,max=100,uuid"`
	CreatorUsername string            `json:"creatorUsername" binding:"required,max=50"`
}

// TenderEdit представляет данные для редактирования существующего тендера
type TenderEdit struct {
	Name        *string            `json:"name" binding:"omitempty,max=100"`
	Description *string            `json:"description" binding:"omitempty,max=500"`
	ServiceType *TenderServiceType `json:"serviceType" binding:"omitempty,oneof=Construction Delivery Manufacture"`
}

func (t *TenderEdit) IsEmpty() bool {
	return t.Name == nil && t.Description == nil && t.ServiceType == nil
}
