package products

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EduardoMark/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	Create(ctx context.Context, sellerID uuid.UUID, name, description string, basePrice float64, auctionEnd time.Time) (uuid.UUID, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*pgstore.Product, error)
	GetAllProducts(ctx context.Context) ([]*pgstore.Product, error)
}

type productService struct {
	pool *pgxpool.Pool
	q    *pgstore.Queries
}

var ErrNotFound = errors.New("not found")

func NewProductService(pool *pgxpool.Pool) Service {
	return &productService{
		pool: pool,
		q:    pgstore.New(pool),
	}
}

func (s *productService) Create(ctx context.Context, sellerID uuid.UUID, name, description string, basePrice float64, auctionEnd time.Time) (uuid.UUID, error) {
	args := pgstore.CreateProductParams{
		SellerID:    sellerID,
		Name:        name,
		Description: description,
		BasePrice:   basePrice,
		AuctionEnd:  auctionEnd,
	}

	id, err := s.q.CreateProduct(ctx, args)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("service.create: %v", err)
	}

	return id, nil
}

func (s *productService) GetProductByID(ctx context.Context, id uuid.UUID) (*pgstore.Product, error) {
	record, err := s.q.GetOneProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("service.getProductByID: %v", err)
	}

	return record, nil
}

func (s *productService) GetAllProducts(ctx context.Context) ([]*pgstore.Product, error) {
	records, err := s.q.GetAllProducts(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("service.getProductByID: %v", err)
	}

	return records, nil
}
