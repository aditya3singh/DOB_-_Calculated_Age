package repository

import (
	"context"

	db "user-api/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

// UserRepository wraps the SQLC Queries to provide a clean data-access interface.
type UserRepository struct {
	q *db.Queries
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{q: q}
}

// Create inserts a new user and returns the created record.
func (r *UserRepository) Create(ctx context.Context, name string, dob pgtype.Date) (db.User, error) {
	return r.q.CreateUser(ctx, db.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

// GetByID fetches a single user by primary key.
func (r *UserRepository) GetByID(ctx context.Context, id int32) (db.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// Update modifies an existing user and returns the updated record.
func (r *UserRepository) Update(ctx context.Context, id int32, name string, dob pgtype.Date) (db.User, error) {
	return r.q.UpdateUser(ctx, db.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  dob,
	})
}

// Delete removes a user by primary key.
func (r *UserRepository) Delete(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}

// List returns a paginated slice of users.
func (r *UserRepository) List(ctx context.Context, limit, offset int32) ([]db.User, error) {
	return r.q.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

// Count returns the total number of users.
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountUsers(ctx)
}
