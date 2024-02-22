package requests

type IsAdminRequest struct {
	UserID int64 `validate:"required"`
}
