package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hunterros-s/algernon/config"
	"github.com/hunterros-s/algernon/server/common"
	"github.com/hunterros-s/algernon/server/listener"
	"github.com/hunterros-s/algernon/server/protocol"
	"github.com/hunterros-s/algernon/server/supervisor"
)

// should not re-create the tcpserver in tcp. just create tcpserver here and add callbacks

type Server struct {
	config     *config.ServerConfig
	listener   *listener.Listener
	signals    chan os.Signal
	supervisor *supervisor.Supervisor
}

func (svr *Server) Start() {
	svr.supervisor.Start()
	svr.listener.Start()
}

func (svr *Server) Wait() {
	<-svr.signals
}

func (svr *Server) Stop() {
	defer svr.config.Logger.Info().Msg("Server shut down.")
	svr.listener.Stop()
	svr.supervisor.Stop()
}

// this server should be the main processing center/thread i think.
func NewServer(cfg *config.ServerConfig) (*Server, error) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	l := listener.NewListener(cfg)

	sv := supervisor.NewSupervisor(cfg.Logger)
	// need to give this access to a central processing channel. it will send packets to that.
	// that will decide what to do to the actual mc server, i.e. change a block, send a chat, leave, join.
	// cant think how it should be structured.
	l.SetOnNewMessage(protocol.GetNewMessageCallback(func(c common.Client, sp common.ServerboundPacket) {
		sv.Handle(common.IncomingEntry{
			Packet: sp,
			Client: c,
		})
	}, cfg.Logger))

	return &Server{
		config:     cfg,
		listener:   l,
		signals:    signals,
		supervisor: sv,
	}, nil
}

// need to create another server that is running. this server will process the packets. do stuff to the server. send packets out to respective clients.
