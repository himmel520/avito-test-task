package httphandler

import "git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"

// errorResponse содержит описание ошибки для ответа в формате JSON
type errorResponse struct {
	Reason string `json:"reason"`
}

// tenderIdURI содержит идентификатор тендера из URI запроса
type tenderIdURI struct {
	ID string `uri:"tenderId" binding:"required,uuid"`
}

// usernameQuery содержит имя пользователя из строки запроса
type usernameQuery struct {
	Username string `form:"username" binding:"required,max=50"`
}

// allTenderFilterQuery содержит параметры для фильтрации и пагинации всех тендеров
type allTenderQuery struct {
	Limit       int32                      `form:"limit,default=5" binding:"omitempty,min=1"`
	Offset      int32                      `form:"offset,default=0" binding:"omitempty,min=0"`
	ServiceType []models.TenderServiceType `form:"service_type" binding:"omitempty,dive,oneof=Construction Delivery Manufacture"`
}

// allTenderFilterQuery содержит параметры для фильтрации и пагинации всех тендеров
type myTenderQuery struct {
	Limit    int32  `form:"limit,default=5" binding:"omitempty,min=1"`
	Offset   int32  `form:"offset,default=0" binding:"omitempty,min=0"`
	Username string `form:"username" binding:"required,max=50"`
}

type editTenderQuery struct {
	Status   models.TenderStatus `form:"status" binding:"required,oneof=Created Published Closed"`
	Username string              `form:"username" binding:"required,max=50"`
}

type rollbackTenderUri struct {
	ID      string `uri:"tenderId" binding:"required,uuid"`
	Version int32  `uri:"version" binding:"required,min=1"`
}
