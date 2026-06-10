package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	db "user-api/db/sqlc"
	"user-api/internal/logger"
	"user-api/internal/models"
	"user-api/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func CalculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()

	birthdayThisYear := time.Date(now.Year(), dob.Month(), dob.Day(), 0, 0, 0, 0, now.Location())
	if now.Before(birthdayThisYear) {
		years--
	}
	return years
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	dob, err := parseDOB(req.DOB)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		logger.Error("failed to create user", zap.Error(err))
		return models.UserResponse{}, fmt.Errorf("could not create user: %w", err)
	}

	logger.Info("user created", zap.Int32("id", user.ID), zap.String("name", user.Name))
	return toUserResponse(user), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int32) (models.UserWithAgeResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserWithAgeResponse{}, ErrUserNotFound
		}
		logger.Error("failed to get user", zap.Int32("id", id), zap.Error(err))
		return models.UserWithAgeResponse{}, fmt.Errorf("could not fetch user: %w", err)
	}

	logger.Info("user fetched", zap.Int32("id", user.ID))
	return toUserWithAgeResponse(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	dob, err := parseDOB(req.DOB)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserResponse{}, ErrUserNotFound
		}
		logger.Error("failed to update user", zap.Int32("id", id), zap.Error(err))
		return models.UserResponse{}, fmt.Errorf("could not update user: %w", err)
	}

	logger.Info("user updated", zap.Int32("id", user.ID))
	return toUserResponse(user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("could not fetch user before delete: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error("failed to delete user", zap.Int32("id", id), zap.Error(err))
		return fmt.Errorf("could not delete user: %w", err)
	}

	logger.Info("user deleted", zap.Int32("id", id))
	return nil
}

func (s *UserService) ListUsers(ctx context.Context, page, limit int) ([]models.UserWithAgeResponse, models.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := s.repo.List(ctx, int32(limit), int32(offset))
	if err != nil {
		logger.Error("failed to list users", zap.Error(err))
		return nil, models.PaginationMeta{}, fmt.Errorf("could not list users: %w", err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		logger.Error("failed to count users", zap.Error(err))
		return nil, models.PaginationMeta{}, fmt.Errorf("could not count users: %w", err)
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	data := make([]models.UserWithAgeResponse, 0, len(users))
	for _, u := range users {
		data = append(data, toUserWithAgeResponse(u))
	}

	logger.Info("users listed", zap.Int("page", page), zap.Int("limit", limit), zap.Int64("total", total))

	return data, models.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

var ErrUserNotFound = errors.New("user not found")

func parseDOB(raw string) (pgtype.Date, error) {
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("invalid dob format, expected YYYY-MM-DD: %w", err)
	}
	return pgtype.Date{Time: t, Valid: true}, nil
}

func toUserResponse(u db.User) models.UserResponse {
	return models.UserResponse{
		ID:   u.ID,
		Name: u.Name,
		DOB:  u.Dob.Time.Format("2006-01-02"),
	}
}

func toUserWithAgeResponse(u db.User) models.UserWithAgeResponse {
	return models.UserWithAgeResponse{
		ID:   u.ID,
		Name: u.Name,
		DOB:  u.Dob.Time.Format("2006-01-02"),
		Age:  CalculateAge(u.Dob.Time),
	}
}
