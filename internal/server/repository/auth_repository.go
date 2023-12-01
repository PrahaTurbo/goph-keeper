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

var ErrAlreadyExist = errors.New("login already exist in database")

type AuthRepository interface {
	SaveUser(ctx context.Context, user models.User) (int, error)
	GetUser(ctx context.Context, login string) (*models.User, error)
}

type authRepo struct {
	pg *pgxpool.Pool
}

func NewAuthRepository(pg *pgxpool.Pool) AuthRepository {
	r := &authRepo{
		pg: pg,
	}

	return r
}

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
