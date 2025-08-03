package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/EduardoMark/gobid/internal/auth/token"
	"github.com/EduardoMark/gobid/internal/jsonutils"
)

type ctxKey string

const UserIDKey ctxKey = "user_id"

func AuthToken(jwtService token.JwtService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const BearerSchema = "Bearer "

			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, BearerSchema) {
				jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
					"error": "unauthorized",
				})
				return
			}

			tokenStr := strings.TrimPrefix(header, BearerSchema)
			claims, err := jwtService.ValidateToken(tokenStr)
			if err != nil {
				jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
					"error": "invalid token",
				})
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
