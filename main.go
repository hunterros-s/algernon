package main

import (
	"github.com/hunterros-s/algernon/config"
	"github.com/hunterros-s/algernon/logger"
	"github.com/hunterros-s/algernon/server"
	handshake "github.com/hunterros-s/algernon/server/network/packet/handshaking"
)

func main() {
	test_packet := &handshake.HandshakePacket{
		ProtocolVersion: 767,
		ServerAddress:   "localhost",
		ServerPort:      25565,
		NextState:       2,
	}

	log := logger.NewLogger()
	cfg := config.NewServerConfig("127.0.0.1", 25565, log)

	server := server.NewServer(cfg)
	server.Start()
}
