package model

type User struct {
	ID           int64
	AppID        int
	Email        string
	PasswordHash []byte
	IsAdmin      bool
}
