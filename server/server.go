package server

import (
	"github.com/hunterros-s/algernon/config"
	packets "github.com/hunterros-s/algernon/server/network/packet"
	"github.com/hunterros-s/algernon/server/network/tcpserver"
)

type Server struct {
	config     *config.ServerConfig
	tcp_server *tcpserver.TCPServer
}

func (svr *Server) Start() {
	svr.tcp_server.Start()
}

func NewServer(cfg *config.ServerConfig) *Server {
	tcp_server := tcpserver.NewTCPServer(cfg, packets.NewPacketHandler())

	return &Server{
		config:     cfg,
		tcp_server: tcp_server,
	}
}
