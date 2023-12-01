package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PrahaTurbo/goph-keeper/internal/server/repository/pg"
)

var ErrNoRows = errors.New("no rows were found")

type SecretRepository interface {
	Create(ctx context.Context, secret *Secret) error
	GetUserSecrets(ctx context.Context, userID int) ([]Secret, error)
	UpdateSecret(ctx context.Context, secret *Secret) error
	DeleteSecret(ctx context.Context, secretID, userID int) error
}

type secretRepo struct {
	pg *pgxpool.Pool
}

func NewSecretRepository(pg *pgxpool.Pool) SecretRepository {
	r := &secretRepo{
		pg: pg,
	}

	return r
}

func (s *secretRepo) Create(ctx context.Context, secret *Secret) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, pg.DefaultQueryTimeout)
	defer cancel()

	stmt := `
INSERT INTO secrets 
    (user_id, 
     type, 
     content, 
     meta_data)
VALUES ($1, $2, $3, $4)
`

	_, err := s.pg.Exec(timeoutCtx, stmt,
		secret.UserID,
		secret.Type,
		secret.Content,
		secret.MetaData)
	if err != nil {
		return err
	}

	return nil
}

func (s *secretRepo) GetUserSecrets(ctx context.Context, userID int) ([]Secret, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, pg.DefaultQueryTimeout)
	defer cancel()

	stmt := `
SELECT id, 
       user_id, 
       type, 
       content,
       meta_data,
       created_at
FROM secrets
WHERE user_id = $1
ORDER BY created_at
`

	rows, err := s.pg.Query(timeoutCtx, stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var secrets []Secret
	for rows.Next() {
		var secret Secret

		err := rows.Scan(
			&secret.ID,
			&secret.UserID,
			&secret.Type,
			&secret.Content,
			&secret.MetaData,
			&secret.CreatedAt)
		if err != nil {
			return nil, err
		}

		secrets = append(secrets, secret)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return secrets, nil
}

func (s *secretRepo) UpdateSecret(ctx context.Context, secret *Secret) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, pg.DefaultQueryTimeout)
	defer cancel()

	stmt := `
UPDATE secrets 
SET type = $1, content = $2, meta_data = $3
WHERE id = $4 AND user_id = $5
`

	tag, err := s.pg.Exec(timeoutCtx, stmt,
		secret.Type,
		secret.Content,
		secret.MetaData,
		secret.ID,
		secret.UserID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNoRows
	}

	return nil
}

func (s *secretRepo) DeleteSecret(ctx context.Context, secretID, userID int) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, pg.DefaultQueryTimeout)
	defer cancel()

	stmt := `
DELETE FROM secrets 
WHERE id = $1 AND user_id = $2
`

	tag, err := s.pg.Exec(timeoutCtx, stmt, secretID, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNoRows
	}

	return nil
}
