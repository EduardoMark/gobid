package api

import (
	"github.com/EduardoMark/gobid/internal/auth"
	"github.com/EduardoMark/gobid/internal/auth/token"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DBPool *pgxpool.Pool
}

func BindRoutes(cfg Config) *chi.Mux {
	r := chi.NewMux()

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Logger)

		setupAuthRoutes(r, cfg.DBPool)
	})

	return r
}

func setupAuthRoutes(r chi.Router, pool *pgxpool.Pool) {
	jwtService := token.NewJwtService()

	authSvc := auth.NewAuthService(pool)
	authHandler := auth.NewAuthHandler(authSvc, jwtService)
	authHandler.RegisterAuthRoutes(r)
}
