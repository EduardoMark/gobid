package auth

import (
	"context"

	"github.com/EduardoMark/gobid/internal/validator"
)

type SignupReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (r *SignupReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(r.Username), "username", "this field cannot be blank")
	eval.CheckField(validator.MinChars(r.Username, 3), "username", "this field must have at least 3 characters")
	eval.CheckField(validator.Matches(r.Email, validator.EmailRX), "email", "this field must be a valid email")
	eval.CheckField(validator.MinChars(r.Password, 8), "password", "this field must have least 8 characters")
	eval.CheckField(
		validator.MinChars(r.Bio, 10) && validator.MaxChars(r.Bio, 255),
		"bio", "this field mus be between 10 and 255 characters",
	)

	return eval
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.Matches(r.Email, validator.EmailRX), "email", "this field must be a valid email")
	eval.CheckField(validator.MinChars(r.Password, 8), "password", "this field must have least 8 characters")

	return eval
}
