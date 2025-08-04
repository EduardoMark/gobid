package users

import (
	"context"
	"time"

	"github.com/EduardoMark/gobid/internal/validator"
	"github.com/google/uuid"
)

type UsersResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
}

func (r *UpdateReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(r.Username), "username", "this field cannot be blank")
	eval.CheckField(validator.MinChars(r.Username, 3), "username", "this field must have at least 3 characters")
	eval.CheckField(validator.Matches(r.Email, validator.EmailRX), "email", "this field must be a valid email")
	eval.CheckField(
		validator.MinChars(r.Bio, 10) && validator.MaxChars(r.Bio, 255),
		"bio", "this field mus be between 10 and 255 characters",
	)

	return eval
}
