package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tangerinefrog/chatter/internal/db"
)

type UsersRepository struct {
	q *db.Queries
}

func NewRepository(pool *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{
		q: db.New(pool),
	}
}

func (r *UsersRepository) Create(ctx context.Context, username, passwordHash string) (*db.User, error) {
	u, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UsersRepository) GetByID(ctx context.Context, id int32) (*db.User, error) {
	u, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &u, nil
}

func (r *UsersRepository) GetByUsername(ctx context.Context, username string) (*db.User, error) {
	u, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &u, nil
}
