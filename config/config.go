package config

import (
	"net"

	"github.com/rs/zerolog"
)

type ServerConfig struct {
	ServerIP   net.IP
	ServerPort int
	TPS        int
	Brand      string
	MOTD       string
	Logger     zerolog.Logger
}

func NewServerConfig(ip_str string, port int, log zerolog.Logger) *ServerConfig {
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
