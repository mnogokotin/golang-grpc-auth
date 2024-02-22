package requests

type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
	AppID    int32  `validate:"required"`
}
