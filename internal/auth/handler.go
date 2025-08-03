package auth

import (
	"errors"
	"net/http"

	"github.com/EduardoMark/gobid/internal/auth/token"
	"github.com/EduardoMark/gobid/internal/jsonutils"
	"github.com/go-chi/chi/v5"
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
