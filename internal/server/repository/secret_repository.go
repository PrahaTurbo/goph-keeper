// Package repository provides an abstraction over users and secrets databases.
package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PrahaTurbo/goph-keeper/internal/server/repository/pg"
)

// ErrNoRows is returned when no rows are found for a query.
var ErrNoRows = errors.New("no rows were found")

// SecretRepository is an interface that defines methods for
// handling secret related operations in the database.
type SecretRepository interface {
	Create(ctx context.Context, secret *Secret) error
	GetUserSecrets(ctx context.Context, userID int) ([]Secret, error)
	UpdateSecret(ctx context.Context, secret *Secret) error
	DeleteSecret(ctx context.Context, secretID, userID int) error
}

type secretRepo struct {
	pg *pgxpool.Pool
}

// NewSecretRepository creates and returns an instance of SecretRepository.
func NewSecretRepository(pg *pgxpool.Pool) SecretRepository {
	r := &secretRepo{
		pg: pg,
	}

	return r
}

// Create implements the Create method of the SecretRepository interface.
// It stores a new secret in the PostgreSQL database.
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

// GetUserSecrets implements the GetUserSecrets method of the SecretRepository interface.
// It retrieves all secrets related to a specific user from the PostgreSQL database.
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

// UpdateSecret implements the UpdateSecret method of the SecretRepository interface.
// It updates an existing secret in the PostgreSQL database.
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

// DeleteSecret implements the DeleteSecret method of the SecretRepository interface.
// It removes a specific secret associated with a User ID from the PostgreSQL database.
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
