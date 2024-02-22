package requests

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
	AppID    int64  `validate:"required"`
}
