package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/EduardoMark/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, username, email, password, bio string) (uuid.UUID, error)
	AuthLogin(ctx context.Context, email, password string) (uuid.UUID, error)
}

type AuthService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewAuthService(pool *pgxpool.Pool) AuthService {
	return AuthService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

var ErrEmailAlreadyExist = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

func (s AuthService) Create(ctx context.Context, username, email, password, bio string) (uuid.UUID, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("create: %v", err)
	}

	args := pgstore.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
		Bio:          bio,
	}

	id, err := s.queries.CreateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.UUID{}, ErrEmailAlreadyExist
		}

		logrus.WithField(
			"err", err.Error(),
		).Error("Failed to create user")

		return uuid.UUID{}, fmt.Errorf("create: %v", err)
	}

	return id, nil
}

func (s AuthService) AuthLogin(ctx context.Context, email, password string) (uuid.UUID, error) {
	record, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, ErrInvalidCredentials
		}

		logrus.WithField(
			"err", err.Error(),
		).Error("AuthLogin")

		return uuid.UUID{}, fmt.Errorf("auth login: %v", err)
	}

	isValidPassword := bcrypt.CompareHashAndPassword([]byte(record.PasswordHash), []byte(password))
	if isValidPassword != nil {
		return uuid.UUID{}, ErrInvalidCredentials
	}

	return record.ID, nil
}
