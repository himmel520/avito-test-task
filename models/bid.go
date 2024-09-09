package models

import "time"

// BidStatus представляет статус предложения
type BidStatus string

const (
    BidStatusCreated  BidStatus = "Created"
    BidStatusPublished BidStatus = "Published"
    BidStatusCanceled BidStatus = "Canceled"
    BidStatusApproved BidStatus = "Approved"
    BidStatusRejected BidStatus = "Rejected"
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
type Bid struct {
    ID          string       `json:"id"`
    Name        string       `json:"name"`
    Description string       `json:"description"`
    Status      BidStatus    `json:"status"`
    TenderID    string       `json:"tenderId"`
    AuthorType  BidAuthorType `json:"authorType"`
    AuthorID    string       `json:"authorId"`
    Version     int          `json:"version"`
    CreatedAt   time.Time    `json:"createdAt"`
}

// BidFeedback представляет отзыв на предложение
type BidFeedback struct {
    Feedback string `json:"feedback"`
}

// BidReview представляет отзыв о предложении
type BidReview struct {
    ID          string    `json:"id"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"createdAt"`
}
