package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/EduardoMark/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Service interface {
	GetOneUser(ctx context.Context, id uuid.UUID) (*pgstore.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, username, email, bio string) (*pgstore.UpdateUserRow, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	pool *pgxpool.Pool
	q    *pgstore.Queries
}

var ErrNotFound = errors.New("not found")
var ErrEmailAlreadyExists = errors.New("email already exists")

func NewUserService(pool *pgxpool.Pool) Service {
	return &userService{
		pool: pool,
		q:    pgstore.New(pool),
	}
}

func (s *userService) GetOneUser(ctx context.Context, id uuid.UUID) (*pgstore.User, error) {
	record, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		logrus.WithField("err", err.Error()).Error("GetOneUser")

		return nil, fmt.Errorf("GetOneUser: %v", err)
	}

	return record, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, username, email, bio string) (*pgstore.UpdateUserRow, error) {
	_, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		logrus.WithField("err", err.Error()).Error("UpdateUser")

		return nil, fmt.Errorf("UpdateUser: %v", err)
	}

	emailInUseByOtherUser, err := s.q.CheckEmailExistsExcludingID(ctx, pgstore.CheckEmailExistsExcludingIDParams{
		ID:    id,
		Email: email,
	})
	if err != nil {
		logrus.WithField("err", err.Error()).Error("UpdateUser - CheckEmailExistsExcludingID")
		return nil, fmt.Errorf("UpdateUser: %v", err)
	}
	if emailInUseByOtherUser {
		return nil, ErrEmailAlreadyExists
	}

	params := pgstore.UpdateUserParams{
		ID:       id,
		Username: username,
		Email:    email,
		Bio:      bio,
	}

	record, err := s.q.UpdateUser(ctx, params)
	if err != nil {
		logrus.WithField("err", err.Error()).Error("UpdateUser")
		return nil, fmt.Errorf("UpdateUser: %v", err)
	}

	return record, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.q.DeleteUser(ctx, id); err != nil {
		logrus.WithField("err", err.Error()).Error("DeleteUser")
		return fmt.Errorf("DeleteUser: %v", err)
	}

	return nil
}
