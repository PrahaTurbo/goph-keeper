// Package repository provides an abstraction over users and secrets databases.
package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PrahaTurbo/goph-keeper/internal/server/models"
	"github.com/PrahaTurbo/goph-keeper/internal/server/repository/pg"
)

const uniqueViolationErrCode = "23505"

// ErrAlreadyExist is returned when a user already exists in the database.
var ErrAlreadyExist = errors.New("login already exist in database")

// AuthRepository is an interface that defines method to
// interact with underlying User related database operations.
type AuthRepository interface {
	SaveUser(ctx context.Context, user models.User) (int, error)
	GetUser(ctx context.Context, login string) (*models.User, error)
}

type authRepo struct {
	pg *pgxpool.Pool
}

// NewAuthRepository creates and returns an instance of AuthRepository.
func NewAuthRepository(pg *pgxpool.Pool) AuthRepository {
	r := &authRepo{
		pg: pg,
	}

	return r
}

// SaveUser implements the SaveUser method of the AuthRepository interface.
// It saves a User record in a PostgreSQL database and handles unique constraint violations.
func (a *authRepo) SaveUser(ctx context.Context, user models.User) (int, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, pg.DefaultQueryTimeout)
	defer cancel()

	stmt := `
INSERT INTO users (login, password)
VALUES ($1, $2)
RETURNING id
`

	var userID int
	err := a.pg.QueryRow(timeoutCtx, stmt, user.Login, user.PasswordHash).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); !ok {
			return 0, err
		}

		if pgErr.Code != uniqueViolationErrCode {
			return 0, err
		}

		return 0, ErrAlreadyExist
	}

	return userID, nil
}

// GetUser implements the GetUser method of the AuthRepository interface.
// It retrieves a User record by login from a PostgreSQL database.
func (a *authRepo) GetUser(ctx context.Context, login string) (*models.User, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, pg.DefaultQueryTimeout)
	defer cancel()

	stmt := `
SELECT id, login, password
FROM users
WHERE login = $1
`

	row := a.pg.QueryRow(timeoutCtx, stmt, login)

	var user models.User
	if err := row.Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		return nil, err
	}

	return &user, nil
}
