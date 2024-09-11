package models

import "time"

// BidStatus представляет статус предложения
type BidStatus string

const (
	BidStatusCreated   BidStatus = "Created"
	BidStatusPublished BidStatus = "Published"
	BidStatusCanceled  BidStatus = "Canceled"
	BidStatusApproved  BidStatus = "Approved"
	BidStatusRejected  BidStatus = "Rejected"
)

// BidDecision представляет решение по предложению
type BidDecision string

const (
	BidDecisionApproved BidDecision = "Approved"
	BidDecisionRejected BidDecision = "Rejected"
)

// BidAuthorType представляет тип автора предложения
type BidAuthorType string

const (
	BidAuthorTypeOrganization BidAuthorType = "Organization"
	BidAuthorTypeUser         BidAuthorType = "User"
)

// Bid представляет информацию о предложении
type BidResponse struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      BidStatus     `json:"status"`
	TenderID    string        `json:"tenderId"`
	AuthorType  BidAuthorType `json:"authorType"`
	AuthorID    string        `json:"authorId"`
	Version     int           `json:"version"`
	CreatedAt   time.Time     `json:"createdAt"`
}

// BidCreate представляет данные для создания нового предложения
type BidCreate struct {
	Name        string        `json:"name" binding:"required,max=100"`
	Description string        `json:"description" binding:"required,max=500"`
	TenderID    string        `json:"tenderId" binding:"required,max=100"`
	AuthorType  BidAuthorType `json:"authorType" binding:"required,oneof=Organization User"`
	AuthorId    string        `json:"authorId" binding:"required,max=100"`
}

// BidEdit представляет данные для редактирования предложения
type BidEdit struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
}

// BidFeedback представляет отзыв на предложение
type BidFeedback string

// BidReview представляет отзыв о предложении
type BidReview struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
