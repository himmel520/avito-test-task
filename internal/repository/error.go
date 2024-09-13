package repository

import "errors"

// DBCode
var (
	// foreign key violation: 23503
	FKViolation = "23503"
	// unique violation: 23505
	UniqueConstraint = "23505"
)

// COMMON
var (
	// ErrUserNotExist указывает на несуществующего пользователя.
	ErrUserNotExist = errors.New("пользователь не существует или некорректен")
	// ErrRelationNotExist указывает на отсутствие прав у пользователя.
	ErrRelationNotExist = errors.New("недостаточно прав для выполнения действия")
)

// TENDER
var (
	// ErrOrganizationDepencyNotFound указывает на отсутствие организации с указанным ID.
	ErrOrganizationDepencyNotFound = errors.New("нельзя создать тендер, так как нет организации с таким id")
	// ErrTenderNotFound указывает на отсутствие тендера.
	ErrTenderNotFound = errors.New("тендер не найден")
	// ErrTenderORVersionNotFound указывает на отсутствие тендера или версии.
	ErrTenderORVersionNotFound = errors.New("тендер или версия не найдены")
	ErrTenderClosed = errors.New("тендер закрыт")
)

// BID
var (
	ErrBidDependencyNotFound = errors.New("нельзя создать предложение, так как нет тендера или пользователя")
	ErrBidUnique = errors.New("на один тендер может быть одно предложение от организации")
	ErrBidTenderNotFound = errors.New("тендер или предложение не найдено")
	ErrBidNotFound = errors.New("предложение не найдено")
	ErrBidORVersionNotFound = errors.New("предложение или версия не найдены")
)

