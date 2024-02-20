package model

type User struct {
	ID           int64
	Email        string
	PasswordHash []byte
	IsAdmin      bool
}
