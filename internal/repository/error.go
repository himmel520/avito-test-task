package repository

import "errors"

// DBCode
var (
	// foreign key violation: 23503
	FKViolation = "23503"
	// unique violation: 23505
	UniqueConstraint = "23505"
)

// TENDER
var (
	ErrUserNotExist = errors.New("пользователь не существует или некорректен")
	ErrRelationNotExist = errors.New("недостаточно прав для выполнения действия")
	ErrPermissionDenied = errors.New("недостаточно прав для выполнения действия")

	ErrOrganizationDepencyNotFound = errors.New("нельзя создать тендер, так как нет организации с таким id")
	ErrTenderNotFound = errors.New("тендер не найден")
)
