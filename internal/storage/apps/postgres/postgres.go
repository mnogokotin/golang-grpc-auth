package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mnogokotin/golang-grpc-auth/internal/domain/model"
	"github.com/mnogokotin/golang-grpc-auth/internal/storage"
	"github.com/mnogokotin/golang-packages/database/postgres"
)

type Postgres struct {
	*postgres.Postgres
}

func (p *Postgres) App(ctx context.Context, appID int64) (model.App, error) {
	const op = "storage.postgres.App"

	var app model.App
	err := p.Db.QueryRow(`SELECT id, name, secret FROM apps 
	WHERE id = $1`, appID).Scan(&app.ID, &app.Name, &app.Secret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return model.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
