package auth

import (
	"errors"
	"net/http"

	"github.com/EduardoMark/gobid/internal/api/middlewares"
	"github.com/EduardoMark/gobid/internal/auth/token"
	"github.com/EduardoMark/gobid/internal/jsonutils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	svc        Service
	jwtService token.JwtService
}

func NewAuthHandler(svc Service, jwtService token.JwtService) AuthHandler {
	return AuthHandler{
		svc:        svc,
		jwtService: jwtService,
	}
}

func (m *AuthHandler) RegisterAuthRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", m.Signup)
		r.Post("/login", m.Login)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthToken(m.jwtService))

			r.Post("/change-password", m.ChangePassword)
		})
	})
}

func (m *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, problems, err := jsonutils.DecodeValidJson[*SignupReq](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := m.svc.Create(
		ctx,
		data.Username,
		data.Email,
		data.Password,
		data.Bio,
	)

	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExist) {
			jsonutils.EncodeJson(w, r, http.StatusConflict, map[string]any{
				"error": "email already exists",
			})
			return
		}

		logrus.WithField("error", err.Error())

		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"id": id,
	})
}

func (m *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, problems, err := jsonutils.DecodeValidJson[*LoginReq](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := m.svc.AuthLogin(ctx, data.Email, data.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
				"error": "invalid credentials",
			})
			return
		}

		logrus.WithField("error", err.Error())

		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	token, err := m.jwtService.GenerateToken(id.String())
	if err != nil {
		logrus.WithField("err", err.Error())

		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"token": token,
	})
}

func (m *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := ctx.Value(middlewares.UserIDKey).(string)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "user ID not found in context",
		})
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid user ID format",
		})
		return
	}

	data, problems, err := jsonutils.DecodeValidJson[*ChangePasswordReq](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	if err := m.svc.ChangePassword(ctx, parsedID, data.CurrentPassword, data.NewPassword); err != nil {
		if errors.Is(err, ErrNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"error": "user not found",
			})
			return
		}

		if errors.Is(err, ErrInvalidCredentials) {
			jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
				"error": "invalid crendentials",
			})
			return
		}

		if errors.Is(err, ErrSamePassword) {
			jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
				"error": "both password is same, send a new password different of the current password",
			})
			return
		}

		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
