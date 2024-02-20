package main

import (
	iapp "github.com/mnogokotin/golang-grpc-auth/internal/app"
	"github.com/mnogokotin/golang-grpc-auth/internal/config"
	"github.com/mnogokotin/golang-grpc-auth/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.New()
	log := logger.New(cfg.Env)

	app := iapp.New(log, cfg.Grpc.Port, cfg.Pg.Url, cfg.TokenTtl)

	go func() {
		app.GrpcServer.MustRun()
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GrpcServer.Stop()
	log.Info("gracefully stopped")
}
