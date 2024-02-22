package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/mnogokotin/golang-grpc-auth/internal/domain/model"
	"github.com/mnogokotin/golang-grpc-auth/internal/storage"
	"github.com/mnogokotin/golang-packages/database/postgres"
	"strings"
)

type Postgres struct {
	*postgres.Postgres
}

func (p *Postgres) SaveUser(ctx context.Context, email string, passwordHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	var user model.User
	err := p.Db.QueryRow(`INSERT INTO users(email, password_hash)
	VALUES($1, $2) RETURNING id`, email, passwordHash).Scan(&user.ID)

	if err != nil {
		var postgresErr pq.PGError
		if errors.As(err, &postgresErr) && strings.Contains(err.Error(), "users_email_key") {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return user.ID, nil
}

func (p *Postgres) User(ctx context.Context, email string) (model.User, error) {
	const op = "storage.postgres.User"

	var user model.User
	err := p.Db.QueryRow(`SELECT id, email, password_hash FROM users 
	WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (p *Postgres) App(ctx context.Context, id int) (model.App, error) {
	const op = "storage.postgres.App"

	var app model.App
	err := p.Db.QueryRow(`SELECT id, name, secret FROM apps 
	WHERE id = $1`, id).Scan(&app.ID, &app.Name, &app.Secret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return model.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (p *Postgres) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	var user model.User
	err := p.Db.QueryRow(`SELECT is_admin FROM users
	WHERE id = $1`, userID).Scan(&user.IsAdmin)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return user.IsAdmin, nil
}
