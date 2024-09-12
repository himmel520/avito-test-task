package repository

import "errors"

// DBCode
var (
	// foreign key violation: 23503
	FKViolation = "23503"
)

// TENDER
var (
	// ErrUserNotExist указывает на несуществующего пользователя.
	ErrUserNotExist = errors.New("пользователь не существует или некорректен")
	// ErrRelationNotExist указывает на отсутствие прав у пользователя.
	ErrRelationNotExist = errors.New("недостаточно прав для выполнения действия")
	
	// ErrOrganizationDepencyNotFound указывает на отсутствие организации с указанным ID.
	ErrOrganizationDepencyNotFound = errors.New("нельзя создать тендер, так как нет организации с таким id")
	// ErrTenderNotFound указывает на отсутствие тендера.
	ErrTenderNotFound = errors.New("тендер не найден")
	// ErrTenderORVersionNotFound указывает на отсутствие тендера или версии.
	ErrTenderORVersionNotFound = errors.New("тендер или версия не найдены")
)
