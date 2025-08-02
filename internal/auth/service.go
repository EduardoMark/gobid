package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/EduardoMark/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, username, email, password, bio string) (uuid.UUID, error)
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
