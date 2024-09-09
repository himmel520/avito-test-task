package models

import "time"

// TenderStatus представляет статус тендера
type TenderStatus string

const (
    TenderStatusCreated  TenderStatus = "Created"
    TenderStatusPublished TenderStatus = "Published"
    TenderStatusClosed   TenderStatus = "Closed"
)

// TenderServiceType представляет вид услуги, к которой относится тендер
type TenderServiceType string

const (
    TenderServiceTypeConstruction TenderServiceType = "Construction"
    TenderServiceTypeDelivery     TenderServiceType = "Delivery"
    TenderServiceTypeManufacture  TenderServiceType = "Manufacture"
)

// Tender представляет информацию о тендере
type Tender struct {
    ID            string             `json:"id"`
    Name          string             `json:"name"`
    Description   string             `json:"description"`
    ServiceType   TenderServiceType  `json:"serviceType"`
    Status        TenderStatus       `json:"status"`
    OrganizationID string            `json:"organizationId"`
    Version       int                `json:"version"`
    CreatedAt     time.Time          `json:"createdAt"`
}

// Username представляет уникальный slug пользователя
type Username string

// OrganizationID представляет уникальный идентификатор организации
type OrganizationID string

