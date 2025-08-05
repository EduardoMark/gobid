package products

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

type ProductHandler struct {
	svc        Service
	jwtService token.JwtService
}

func NewProductHandler(svc Service, jwt token.JwtService) ProductHandler {
	return ProductHandler{
		svc:        svc,
		jwtService: jwt,
	}
}

func (m *ProductHandler) RegisterProductsRoutes(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthToken(m.jwtService))

			r.Post("/", m.Create)
			r.Get("/{id}", m.GetOne)
			r.Get("/", m.GetAll)
		})
	})
}

func (m *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := ctx.Value(middlewares.UserIDKey).(string)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "user ID not found in context",
		})
		return
	}

	sellerID, err := uuid.Parse(id)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid user ID format",
		})
		return
	}

	data, problems, err := jsonutils.DecodeValidJson[*CreateProductReq](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	productID, err := m.svc.Create(
		ctx,
		sellerID,
		data.Name,
		data.Description,
		data.BasePrice,
		data.AuctionEnd,
	)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to create product auction try again later",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"id": productID,
	})
}

func (m *ProductHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid user ID format",
		})
		return
	}

	record, err := m.svc.GetProductByID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"error": "product not found",
			})
			return
		}

		logrus.WithField("err", err.Error()).Error("Handler.GetOne")

		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	res := ProductResponse{
		ID:          record.ID,
		SellerID:    record.SellerID,
		Name:        record.Name,
		Description: record.Description,
		BasePrice:   record.BasePrice,
		AuctionEnd:  record.AuctionEnd,
		IsSold:      record.IsSold,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"product": res,
	})
}

func (m *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	records, err := m.svc.GetAllProducts(ctx)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"error": "product not found",
			})
			return
		}

		logrus.WithField("err", err.Error()).Error("Handler.GetOne")

		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	res := make([]ProductResponse, len(records))
	for i, record := range records {
		res[i] = ProductResponse{
			ID:          record.ID,
			SellerID:    record.SellerID,
			Name:        record.Name,
			Description: record.Description,
			BasePrice:   record.BasePrice,
			AuctionEnd:  record.AuctionEnd,
			IsSold:      record.IsSold,
			CreatedAt:   record.CreatedAt,
			UpdatedAt:   record.UpdatedAt,
		}
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"products": res,
	})
}
