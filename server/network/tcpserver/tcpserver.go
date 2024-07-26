package tcpserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/hunterros-s/algernon/config"
	"github.com/hunterros-s/algernon/logger"
	"github.com/hunterros-s/algernon/server/network/session"
)

// Handler interface for handling connections
type Handler interface {
	Handle(ctx context.Context, conn net.Conn, log logger.Logger)
}

// PacketHandler interface for handling packets
type PacketHandler interface {
	ReadPacket(ctx context.Context, session *session.Session, byteslice []byte, log logger.Logger)
}

// TCPServer represents a TCP server
type TCPServer struct {
	listener net.Listener
	config   *config.ServerConfig
	handler  Handler
	wg       sync.WaitGroup
}

func NewTCPServer(config *config.ServerConfig, packetHandler PacketHandler) (*TCPServer, error) {
	listener, err := net.Listen("tcp4", ":"+strconv.Itoa(config.ServerPort))
	if err != nil {
		return nil, fmt.Errorf("error creating new server listener: %w", err)
	}

	return &TCPServer{
		listener: listener,
		config:   config,
		handler:  NewConnectionHandler(packetHandler),
	}, nil
}

// Ensure that tcpserver.Start() handles graceful shutdowns effectively.
// Start begins accepting connections
func (s *TCPServer) Start(ctx context.Context) error {
	s.config.Logger.Info().Msgf("Accepting connections at %s:%d", s.config.ServerIP, s.config.ServerPort)

	go func() {
		<-ctx.Done()
		s.listener.Close()
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			s.config.Logger.Error().Err(err).Msg("Accept error")
			continue
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.handler.Handle(ctx, conn, s.config.Logger)
		}()
	}
}

func (s *TCPServer) Stop() {
	s.listener.Close()
	s.Wait() // Ensure all connections have finished
}

// Wait waits for all connections to close
func (s *TCPServer) Wait() {
	s.wg.Wait()
}

// connectionHandler handles individual connections
type connectionHandler struct {
	packetHandler PacketHandler
	bufferPool    sync.Pool
}

// NewConnectionHandler creates a new connectionHandler
func NewConnectionHandler(packetHandler PacketHandler) *connectionHandler {
	return &connectionHandler{
		packetHandler: packetHandler,
		bufferPool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 4096)
			},
		},
	}
}

// Handle handles an individual connection
func (h *connectionHandler) Handle(ctx context.Context, conn net.Conn, log logger.Logger) {
	defer conn.Close()
	log.Info().Msgf("Serving %s", conn.RemoteAddr().String())

	buffer := h.bufferPool.Get().([]byte)
	defer h.bufferPool.Put(buffer)

	s := session.NewSession(conn)

	for {
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Minute)); err != nil {
			log.Error().Err(err).Msg("Failed to set read deadline")
			return
		}

		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Info().Msgf("Connection closed by client: %s", conn.RemoteAddr().String())
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Info().Msgf("Connection timeout: %s", conn.RemoteAddr().String())
			} else {
				log.Error().Err(err).Msgf("Error reading from connection: %s", conn.RemoteAddr().String())
			}
			return
		}

		subLogger := log.With().Str("addr", conn.RemoteAddr().String()).Logger()
		h.packetHandler.ReadPacket(ctx, s, buffer[:n], &subLogger)

		select {
		case <-ctx.Done():
			log.Info().Msgf("Server closing connection with: %s", conn.RemoteAddr().String())
			return
		default:
		}
	}
}
