package server

import (
	"context"
	"fmt"

	"github.com/hunterros-s/algernon/config"
	packets "github.com/hunterros-s/algernon/server/network/packet"
	"github.com/hunterros-s/algernon/server/network/tcpserver"
)

type Server struct {
	config    *config.ServerConfig
	tcpServer *tcpserver.TCPServer
}

func (svr *Server) Start(context context.Context) {
	svr.tcpServer.Start(context)
}

func (svr *Server) Stop() {
	svr.config.Logger.Info().Msg("Shutting down server...")
	svr.tcpServer.Stop()
}

func NewServer(cfg *config.ServerConfig) (*Server, error) {
	tcpServer, err := tcpserver.NewTCPServer(cfg, packets.NewPacketHandler())
	if err != nil {
		return nil, fmt.Errorf("error creating tcp server:%v", err.Error())
	}

	return &Server{
		config:    cfg,
		tcpServer: tcpServer,
	}, nil
}
