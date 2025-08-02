package auth

import (
	"errors"
	"net/http"

	"github.com/EduardoMark/gobid/internal/jsonutils"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	svc Service
}

func NewAuthHandler(svc Service) AuthHandler {
	return AuthHandler{
		svc: svc,
	}
}

func (m *AuthHandler) RegisterAuthRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", m.Signup)
	})
}

func (m *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[*SignupReq](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := m.svc.Create(
		r.Context(),
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

		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"id": id,
	})
}
