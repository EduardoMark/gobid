package users

import (
	"errors"
	"net/http"

	"github.com/EduardoMark/gobid/internal/api/middlewares"
	"github.com/EduardoMark/gobid/internal/auth/token"
	"github.com/EduardoMark/gobid/internal/jsonutils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	s          Service
	jwtService token.JwtService
}

func NewUserHandler(s Service, jwtService token.JwtService) UserHandler {
	return UserHandler{
		s:          s,
		jwtService: jwtService,
	}
}

func (m *UserHandler) RegisterUserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthToken(m.jwtService))

			r.Get("/{id}", m.GetOne)
			r.Put("/{id}", m.Update)
			r.Delete("/{id}", m.Delete)
		})
	})
}

func (m *UserHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid uuid type",
		})
		return
	}

	record, err := m.s.GetOneUser(ctx, parsedID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"error": "user not found",
			})
			return
		}

		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	res := UsersResponse{
		ID:        record.ID,
		Username:  record.Username,
		Email:     record.Email,
		Bio:       record.Bio,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"user": res,
	})
}

func (m *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid uuid type",
		})
		return
	}

	data, problems, err := jsonutils.DecodeValidJson[*UpdateReq](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	updatedRecord, err := m.s.UpdateUser(
		ctx,
		parsedID,
		data.Username,
		data.Email,
		data.Bio,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"error": "user not found",
			})
			return
		}

		if errors.Is(err, ErrEmailAlreadyExists) {
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

	res := UsersResponse{
		ID:        updatedRecord.ID,
		Username:  updatedRecord.Username,
		Email:     updatedRecord.Email,
		Bio:       updatedRecord.Bio,
		CreatedAt: updatedRecord.CreatedAt,
		UpdatedAt: updatedRecord.UpdatedAt,
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"user": res,
	})
}

func (m *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid uuid type",
		})
		return
	}

	if err := m.s.DeleteUser(ctx, parsedID); err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
