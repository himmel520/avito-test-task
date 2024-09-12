package repository

import "errors"

// DBCode
var (
	// foreign key violation: 23503
	FKViolation = "23503"
)

// TENDER
var (
	ErrUserNotExist = errors.New("пользователь не существует или некорректен")
	ErrRelationNotExist = errors.New("недостаточно прав для выполнения действия")

	ErrOrganizationDepencyNotFound = errors.New("нельзя создать тендер, так как нет организации с таким id")
	ErrTenderNotFound = errors.New("тендер не найден")
	ErrTenderORVersionNotFound = errors.New("тендер или версия не найдены")
)
