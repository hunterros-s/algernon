package listener

import (
	"fmt"
	"io"
	"sync"

	"github.com/hunterros-s/algernon/config"
	"github.com/hunterros-s/algernon/server/common"
	"github.com/hunterros-s/algernon/tcpserver"
	"github.com/rs/zerolog"
)

type Listener struct {
	tcp     *tcpserver.TCPServer
	clients sync.Map
	logger  zerolog.Logger

	// Callbacks
	onNewClient          func(c common.Client)
	onClientDisconnected func(c common.Client, err error)
	onClientError        func(c common.Client, err error)
	onNewMessage         func(c common.Client, b []byte)
	onListenerStart      func(s *Listener)
	onListenerStop       func(s *Listener)
}

func (l *Listener) StoreClient(c *tcpserver.Client, client *client) {
	l.clients.Store(c.GetUUID(), client)
}

func (l *Listener) LoadClient(c *tcpserver.Client) (*client, bool) {
	val, ok := l.clients.Load(c.GetUUID())
	if !ok {
		return nil, false
	}
	client, ok := val.(*client)
	return client, ok
}

func (l *Listener) DeleteClient(c *tcpserver.Client) {
	l.clients.Delete(c.GetUUID())
}

// SetOnNewClient sets the callback for when a new client connects.
func (l *Listener) SetOnNewClient(callback func(c common.Client)) {
	l.onNewClient = callback
}

// SetOnClientDisconnected sets the callback for when a client disconnects.
func (l *Listener) SetOnClientDisconnected(callback func(c common.Client, err error)) {
	l.onClientDisconnected = callback
}

// SetOnClientError sets the callback for when a client encounters an error.
func (l *Listener) SetOnClientError(callback func(c common.Client, err error)) {
	l.onClientError = callback
}

// SetOnNewMessage sets the callback for when a new message is received.
func (l *Listener) SetOnNewMessage(callback func(c common.Client, b []byte)) {
	l.onNewMessage = callback
}

// SetOnListenerStart sets the callback for when the listener starts.
func (l *Listener) SetOnListenerStart(callback func(s *Listener)) {
	l.onListenerStart = callback
}

// SetOnListenerStop sets the callback for when the listener stops.
func (l *Listener) SetOnListenerStop(callback func(s *Listener)) {
	l.onListenerStop = callback
}

// TCP Callbacks

func (l *Listener) newClientTCP(c *tcpserver.Client) {
	client := newClient(c, l.logger)
	l.StoreClient(c, client)
	client.Logger.Debug().Msg("Client connected")
	if l.onNewClient != nil {
		l.onNewClient(client)
	}
}

func (l *Listener) clientClosedTCP(c *tcpserver.Client, err error) {
	client, ok := l.LoadClient(c)
	if !ok {
		return
	}
	if err != nil && err != io.EOF {
		client.Logger.Error().Err(err).Msg("Client disconnected with error")
		if l.onClientError != nil {
			l.onClientError(client, err)
		}
	} else {
		client.Logger.Debug().Msg("Client disconnected")
		if l.onClientDisconnected != nil {
			l.onClientDisconnected(client, err)
		}
	}

	l.DeleteClient(c)
}

func (l *Listener) onNewMessageTCP(c *tcpserver.Client, b []byte) {
	client, ok := l.LoadClient(c)
	if !ok {
		return
	}

	if l.onNewMessage != nil {
		l.onNewMessage(client, b)
	}
}

func (l *Listener) onListenerStartTCP(s *tcpserver.TCPServer) {
	l.logger.Info().Msg("Listener started")
	if l.onListenerStart != nil {
		l.onListenerStart(l)
	}
}

func (l *Listener) onListenerStopTCP(s *tcpserver.TCPServer) {
	l.logger.Info().Msg("Listener stopped")
	if l.onListenerStop != nil {
		l.onListenerStop(l)
	}
}

// packet handler should be passed in as a parameter here. should be set as a function that the new message function calls once it converts.
func NewListener(cfg *config.ServerConfig) *Listener {

	address := fmt.Sprintf("%s:%d", cfg.ServerIP.String(), cfg.ServerPort)

	tcpsvr := tcpserver.NewServer(address)
	listener := &Listener{
		tcp:    tcpsvr,
		logger: cfg.Logger.With().Str("server address", address).Logger(),
	}

	// setup tcpserver callbacks
	tcpsvr.SetOnNewClient(listener.newClientTCP)
	tcpsvr.SetOnClientClosed(listener.clientClosedTCP)
	tcpsvr.SetOnNewMessage(listener.onNewMessageTCP)
	tcpsvr.SetOnServerStart(listener.onListenerStartTCP)
	tcpsvr.SetOnServerStop(listener.onListenerStopTCP)

	return listener
}

func (l *Listener) Start() {
	l.tcp.Start()
}

func (l *Listener) Stop() {
	l.tcp.Stop()
}
