package repository

import (
	"context"

	db "user-api/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{q: q}
}

func (r *UserRepository) Create(ctx context.Context, name string, dob pgtype.Date) (db.User, error) {
	return r.q.CreateUser(ctx, db.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

func (r *UserRepository) GetByID(ctx context.Context, id int32) (db.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *UserRepository) Update(ctx context.Context, id int32, name string, dob pgtype.Date) (db.User, error) {
	return r.q.UpdateUser(ctx, db.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  dob,
	})
}

func (r *UserRepository) Delete(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}

func (r *UserRepository) List(ctx context.Context, limit, offset int32) ([]db.User, error) {
	return r.q.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountUsers(ctx)
}
