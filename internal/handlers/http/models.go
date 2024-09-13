package httphandler

import "git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"

// COMMON

// errorResponse содержит описание ошибки для ответа в формате JSON
type errorResponse struct {
	Reason string `json:"reason"`
}

type PaginationQuery struct {
	Limit  int32 `form:"limit,default=5" binding:"omitempty,min=1"`
	Offset int32 `form:"offset,default=0" binding:"omitempty,min=0"`
}

// type myQuery struct {
// 	Limit    int32  `form:"limit,default=5" binding:"omitempty,min=1"`
// 	Offset   int32  `form:"offset,default=0" binding:"omitempty,min=0"`
// 	Username string `form:"username" binding:"required,max=50"`
// }

// usernameQuery содержит имя пользователя из строки запроса
type UsernameQuery struct {
	Username string `form:"username" binding:"required,max=50"`
}

type myQuery struct {
	PaginationQuery
	UsernameQuery
}

// ----------------------------------------------------------------------------

// TENDER

// tenderIdURI содержит идентификатор тендера из URI запроса
type tenderIdURI struct {
	ID string `uri:"tenderId" binding:"required,uuid"`
}

// allTenderFilterQuery содержит параметры для фильтрации и пагинации всех тендеров
type allTenderQuery struct {
	Limit       int32                      `form:"limit,default=5" binding:"omitempty,min=1"`
	Offset      int32                      `form:"offset,default=0" binding:"omitempty,min=0"`
	ServiceType []models.TenderServiceType `form:"service_type" binding:"omitempty,dive,oneof=Construction Delivery Manufacture"`
}

// editTenderQuery содержит параметры для редактирования статуса тендера
type updateTenderStatusQuery struct {
	Status   models.TenderStatus `form:"status" binding:"required,oneof=Created Published Closed"`
	Username string              `form:"username" binding:"required,max=50"`
}

// rollbackTenderUri содержит параметры для отката тендера к определенной версии
type rollbackTenderUri struct {
	ID      string `uri:"tenderId" binding:"required,uuid"`
	Version int32  `uri:"version" binding:"required,min=1"`
}

// ----------------------------------------------------------------------------

// BID

type bidTenderIdURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type bidIdURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type updateBidStatusQuery struct {
	Status   models.BidStatus `form:"status" binding:"required,oneof=Created Published Closed"`
	Username string           `form:"username" binding:"required,max=50"`
}

// rollbackTenderUri содержит параметры для отката тендера к определенной версии
type rollbackBidUri struct {
	ID      string `uri:"id" binding:"required,uuid"`
	Version int32  `uri:"version" binding:"required,min=1"`
}

// ----------------------------------------------------------------------------

// FEEDBACK

type feedbackQuery struct {
	BidFeedback models.BidFeedback `form:"bidFeedback" binding:"required,max=500"`
	Username    string             `form:"username" binding:"required,max=50"`
}
