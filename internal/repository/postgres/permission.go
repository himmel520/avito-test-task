package postgres

import (
	"context"
	"errors"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/jackc/pgx/v5"
)

// СheckOrganizationPermission проверяет, имеет ли пользователь права доступа к организации
func (p *Postgres) СheckOrganizationPermission(ctx context.Context, organizationID *models.OrganizationID, username string) error {
	var existsRelation bool

	err := p.DB.QueryRow(ctx, `
	SELECT
		EXISTS (
			SELECT 1
			FROM organization_responsible
			WHERE user_id = e.id
			AND organization_id = $2
		) AS exists_relation
	FROM employee e
	WHERE username = $1;`, username, organizationID).Scan(&existsRelation)

	switch {
	// пользователь не существует или некорректен.
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	// нет связи пользователь и организация
	case !existsRelation:
		return repository.ErrRelationNotExist
	}

	return err
}

// IsUserResponsible проверяет, имеет ли пользователь права доступа к любой организации
func (p *Postgres) IsUserResponsible(ctx context.Context, userId string) error {
	var existsRelation bool

	err := p.DB.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1
            FROM organization_responsible orr
            WHERE orr.user_id = $1
        ) AS exists_relation
        FROM employee e
        WHERE e.id = $1`, userId).Scan(&existsRelation)
	switch {
	// пользователь не существует или некорректен.
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	// нет связи пользователь и организация
	case !existsRelation:
		return repository.ErrRelationNotExist
	}

	return err
}

// IsTenderCreator проверяет, является ли пользователь создателем указанного тендера по его имени
func (p *Postgres) IsTenderCreatorByName(ctx context.Context, tenderId, creatorUsername string) error {
	var isCreator bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
		SELECT 1
		FROM tender
			WHERE creator_username = $1 AND id = $2
	) AS is_creator
	FROM employee
		WHERE username = $1;`, creatorUsername, tenderId).Scan(&isCreator)

	switch {
	// пользователь не существует или некорректен
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	// нет связи пользователь и тендер
	case !isCreator:
		return repository.ErrRelationNotExist
	}

	return err
}

// IsTenderCreator проверяет, является ли пользователь создателем указанного тендера по его id
func (p *Postgres) IsTenderCreatorByID(ctx context.Context, tenderId, creatorId string) error {
	var isCreator bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
        SELECT 1
        FROM tender t
        WHERE t.creator_username = e.username
        AND t.id = $1
    ) AS is_creator
    FROM employee e
    	WHERE e.id = $2;`, tenderId, creatorId).Scan(&isCreator)

	switch {
	// пользователь не существует или некорректен
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	// нет связи пользователь и тендер
	case !isCreator:
		return repository.ErrRelationNotExist
	}

	return err
}

// IsBidCreatorByName проверяет, является ли пользователь создателем указанного предложения по имени
func (p *Postgres) IsBidCreatorByName(ctx context.Context, bidID, creatorUsername string) error {
	var isCreator bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
		SELECT 1
		FROM bid b 
			WHERE author_id = e.id and id = $2
	) AS is_creator
	FROM employee e
		WHERE username = $1`, creatorUsername, bidID).Scan(&isCreator)

	switch {
	// пользователь не существует или некорректен
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	// нет связи пользователь и предложение
	case !isCreator:
		return repository.ErrRelationNotExist
	}

	return err
}

// GetUserIDByName возвращает id пользователя по его имени
func (p *Postgres) GetUserIDByName(ctx context.Context, username string) (string, error) {
	var userId string
	err := p.DB.QueryRow(ctx, `
	select 
		id 
	from employee e 
	where e.username = $1;
	`, username).Scan(&userId)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", repository.ErrUserNotExist
	}

	return userId, err
}

// IsUserResponsibleForTender проверяет, имеет ли пользователь права на организацию тендера
func (p *Postgres) IsUserResponsibleForTender(ctx context.Context, tenderID, username string) error {
	var isRelated bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
		SELECT 1
		FROM tender t
		JOIN organization_responsible orr ON orr.user_id = e.id
		WHERE t.id = $2
		AND orr.organization_id = t.organization_id
	) AS is_related
	from employee e
		where username = $1; `, username, tenderID).Scan(&isRelated)

	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrUserNotExist
	}

	if !isRelated {
		return repository.ErrRelationNotExist
	}

	return err
}

// IsUserResponsibleForTenderByBidID проверяет, имеет ли пользователь права на организацию тендера, к которому относится предложение
func (p *Postgres) IsUserResponsibleForTenderByBidID(ctx context.Context, bidID, username string) error {
	var isRelated bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
		SELECT 1
			FROM tender t
			JOIN organization_responsible orr ON orr.organization_id = t.organization_id
			JOIN employee e ON e.id = orr.user_id
			WHERE t.id = (
				SELECT tender_id
				FROM bid
				WHERE id = $2
		)
		AND e.username = $1
	) AS is_related;`, username, bidID).Scan(&isRelated)

	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrUserNotExist
	}

	if !isRelated {
		return repository.ErrRelationNotExist
	}

	return err
}

// IsUserResponsibleForAuthorBid проверяет, относится ли пользователь к организации автора предложения
func (p *Postgres) IsUserResponsibleForAuthorBid(ctx context.Context, bidID, username string) error {
	var isRelated bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
		SELECT 1
			FROM bid b
			JOIN organization_responsible orr ON orr.user_id = e.id
			JOIN organization_responsible o ON o.organization_id = orr.organization_id
				WHERE b.id = $2
				AND o.user_id = b.author_id
	) AS is_related
	FROM employee e
		WHERE e.username = $1;`, username, bidID).Scan(&isRelated)

	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrUserNotExist
	}

	if !isRelated {
		return repository.ErrRelationNotExist
	}

	return err
}
