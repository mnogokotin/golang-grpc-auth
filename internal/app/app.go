package app

import (
	grpcapp "github.com/mnogokotin/golang-grpc-auth/internal/app/grpc"
	"github.com/mnogokotin/golang-grpc-auth/internal/services/auth"
	"github.com/mnogokotin/golang-grpc-auth/internal/storage/postgres"
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
	storage := &postgres.Postgres{
		Postgres: ppg,
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GrpcServer: grpcApp,
	}
}
