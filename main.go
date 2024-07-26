package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hunterros-s/algernon/config"
	"github.com/hunterros-s/algernon/logger"
	"github.com/hunterros-s/algernon/server"
)

func main() {
	log := logger.NewLogger()
	cfg := config.NewServerConfig("127.0.0.1", 25565, log)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := server.NewServer(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Unable to start server")
		os.Exit(1)
	}

	// Create a channel to listen for OS signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-signals
		cancel()
		server.Stop()
	}()

	server.Start(ctx)
}
