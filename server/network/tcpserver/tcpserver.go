package tcpserver

import (
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/hunterros-s/algernon/config"
	"github.com/hunterros-s/algernon/logger"
)

// Handler interface for handling connections
type Handler interface {
	Handle(conn net.Conn, log logger.Logger)
}

type TCPServer struct {
	listener net.Listener
	config   *config.ServerConfig

	handler Handler
}

func NewTCPServer(config *config.ServerConfig, packetHandler PacketHandler) *TCPServer {
	listener, err := net.Listen("tcp4", ":"+strconv.Itoa(config.ServerPort))
	if err != nil {
		config.Logger.Error().Err(err).Msg("error creating new server listener")
		return nil
	}

	return &TCPServer{
		listener: listener,
		config:   config,
		handler:  NewConnectionHandler(packetHandler),
	}
}

func (s *TCPServer) Start() {
	s.config.Logger.Info().Msgf("Accepting connections at %s:%s", s.config.ServerIP.String(), strconv.Itoa(s.config.ServerPort))
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.config.Logger.Error().Msg(fmt.Sprintf("Accept error: %v", err))
			continue
		}
		go s.handler.Handle(conn, s.config.Logger)
	}
}

type PacketHandler interface {
	ReadPacket(conn net.Conn, byteslice []byte, log logger.Logger)
}

type handle_conn struct {
	packetHandler PacketHandler
}

func (h *handle_conn) Handle(conn net.Conn, log logger.Logger) {
	log.Info().Msgf("Serving %s", conn.RemoteAddr().String())
	defer conn.Close()

	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Info().Msgf("Connection closed by client: %s", conn.RemoteAddr().String())
			} else {
				log.Error().Err(err).Msgf("Error reading from connection: %s", conn.RemoteAddr().String())
			}
			break
		}
		sublogger := log.With().Str("addr", conn.RemoteAddr().String()).Logger()
		h.packetHandler.ReadPacket(conn, buffer[:n], &sublogger)
	}
}

func NewConnectionHandler(packetHandler PacketHandler) *handle_conn {
	return &handle_conn{
		packetHandler: packetHandler,
	}
}
