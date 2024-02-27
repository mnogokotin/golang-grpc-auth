package app

import (
	grpcapp "github.com/mnogokotin/golang-grpc-auth/internal/app/grpc"
	"github.com/mnogokotin/golang-grpc-auth/internal/services/auth"
	apostgres "github.com/mnogokotin/golang-grpc-auth/internal/storage/apps/postgres"
	upostgres "github.com/mnogokotin/golang-grpc-auth/internal/storage/users/postgres"
	ppostgres "github.com/mnogokotin/golang-packages/database/postgres"
	"log/slog"
	"time"
)

type App struct {
	GrpcServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	pgUrl string,
	tokenTTL time.Duration,
) *App {
	ppg, err := ppostgres.New(pgUrl)
	if err != nil {
		panic(err)
	}
	userStorage := &upostgres.Postgres{
		Postgres: ppg,
	}
	appStorage := &apostgres.Postgres{
		Postgres: ppg,
	}

	authService := auth.New(log, userStorage, userStorage, appStorage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GrpcServer: grpcApp,
	}
}
