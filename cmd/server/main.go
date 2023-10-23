package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pvarentsov/powtcp/internal/app/server"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/log"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/tcp"
	"github.com/pvarentsov/powtcp/internal/pkg/service"
)

func main() {
	op := "main"
	ctx, cancel := context.WithCancel(context.Background())

	logger := log.New(log.Opts{
		Level: log.LevelDebug,
		Json:  false,
	})

	service := service.NewServer(service.ServerOpts{
		Logger:       logger,
		ErrorChecker: tcp.NewConnErrorChecker(),
	})

	server, err := server.Listen(ctx, server.Opts{
		Address: ":8080",
		Logger:  logger,
		Service: service,
	})
	if err != nil {
		logger.Error(err.Error(), "op", op)
		os.Exit(1)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	<-signalChannel
	cancel()
	server.Shutdown()
}
