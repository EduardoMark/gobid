package products

import (
	"context"
	"time"

	"github.com/EduardoMark/gobid/internal/validator"
	"github.com/google/uuid"
)

type CreateProductReq struct {
	SellerID    string    `json:"seller_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	AuctionEnd  time.Time `json:"auction_end"`
}

const minAuctionDuration = time.Hour * 2

func (r *CreateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(r.Name), "name", "this field cannot be blank")
	eval.CheckField(validator.NotBlank(r.Description), "description", "this field cannot be blank")
	eval.CheckField(
		validator.MinChars(r.Description, 10) && validator.MaxChars(r.Description, 255),
		"description", "this field must have a length between 10 and 255",
	)
	eval.CheckField(r.BasePrice > 0, "base_price", "this field grather than 0")
	eval.CheckField(time.Until(r.AuctionEnd) >= minAuctionDuration, "auction_end", "must be at least two hours duration")

	return eval
}

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	SellerID    uuid.UUID `json:"seller_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	AuctionEnd  time.Time `json:"auction_end"`
	IsSold      bool      `json:"is_sold"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
