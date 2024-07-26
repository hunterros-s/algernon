package config

import (
	"net"

	"github.com/hunterros-s/algernon/logger"
)

type ServerConfig struct {
	ServerIP   net.IP
	ServerPort int
	TPS        int
	Brand      string
	MOTD       string
	Logger     logger.Logger
}

func NewServerConfig(ip_str string, port int, log logger.Logger) *ServerConfig {
	defer log.Info().Msg("Created server config.")

	ip := net.ParseIP(ip_str)
	if ip == nil {
		log.Error().Msgf("Invalid IP address: %s", ip_str)
	}

	return &ServerConfig{
		ServerIP:   ip,
		ServerPort: port,
		TPS:        20,
		Brand:      "algernon",
		MOTD:       "algernon dev server",
		Logger:     log,
	}
}
