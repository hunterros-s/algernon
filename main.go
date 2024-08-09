package main

import (
	"os"

	"github.com/hunterros-s/algernon/config"
	"github.com/hunterros-s/algernon/logger"
	"github.com/hunterros-s/algernon/server"
)

func main() {
	log := logger.NewLogger()
	cfg := config.NewServerConfig("127.0.0.1", 25565, log)

	server, err := server.NewServer(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Unable to start server")
		os.Exit(1)
	}

	server.Start()
	server.Wait()
	server.Stop()
}
