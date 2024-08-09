package tcpserver

import (
	"net"
	"sync"

	"github.com/google/uuid"
)

type Client struct {
	conn   net.Conn
	server *TCPServer
	uuid   uuid.UUID
	send   chan []byte
}

func (c *Client) Send(message []byte) {
	c.send <- message
}

func (c *Client) GetUUID() uuid.UUID {
	return c.uuid
}

func (c *Client) GetIP() string {
	return c.conn.RemoteAddr().String()
}

func NewClient(conn net.Conn, server *TCPServer) *Client {
	return &Client{
		conn:   conn,
		server: server,
		uuid:   uuid.New(),
		send:   make(chan []byte),
	}
}

type TCPServer struct {
	address  string
	listener net.Listener
	clients  map[uuid.UUID]*Client
	mutex    sync.Mutex
	wg       sync.WaitGroup
	stopChan chan struct{}

	// Callbacks
	onNewClient    func(*Client)         // callback for when a new client connects
	onClientClosed func(*Client, error)  // callback for when a client disconnects
	onNewMessage   func(*Client, []byte) // callback for when a new message is received
	onServerStart  func(*TCPServer)      // callback for when the server starts
	onServerStop   func(*TCPServer)      // callback for when the server stops
}

func NewServer(address string) *TCPServer {
	return &TCPServer{
		address:  address,
		clients:  make(map[uuid.UUID]*Client),
		stopChan: make(chan struct{}),
	}
}

// SetOnNewClient sets the callback for when a new client connects
func (s *TCPServer) SetOnNewClient(callback func(*Client)) {
	s.onNewClient = callback
}

// SetOnClientClosed sets the callback for when a client disconnects
func (s *TCPServer) SetOnClientClosed(callback func(*Client, error)) {
	s.onClientClosed = callback
}

// SetOnNewMessage sets the callback for when a new message is received
func (s *TCPServer) SetOnNewMessage(callback func(*Client, []byte)) {
	s.onNewMessage = callback
}

// SetOnServerStart sets the callback for when the server starts
func (s *TCPServer) SetOnServerStart(callback func(*TCPServer)) {
	s.onServerStart = callback
}

// SetOnServerStop sets the callback for when the server stops
func (s *TCPServer) SetOnServerStop(callback func(*TCPServer)) {
	s.onServerStop = callback
}

func (s *TCPServer) GetAddress() string {
	return s.address
}

func (s *TCPServer) IsStopped() bool {
	select {
	case <-s.stopChan:
		return true
	default:
		return false
	}

}

func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp4", s.address)
	if err != nil {
		return err
	}
	s.listener = listener

	if s.onServerStart != nil {
		s.onServerStart(s)
	}

	s.wg.Add(1)
	go s.acceptConnections()

	return nil
}

func (s *TCPServer) Stop() {
	close(s.stopChan)
	s.listener.Close()

	s.mutex.Lock()
	for _, client := range s.clients {
		client.conn.Close()
		close(client.send)
	}
	s.mutex.Unlock()

	s.wg.Wait()

	if s.onServerStop != nil {
		s.onServerStop(s)
	}
}

func (s *TCPServer) acceptConnections() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if s.IsStopped() {
			return
		}
		if err != nil {
			// should make a callback for an error here
			continue
		}

		client := NewClient(conn, s)

		s.mutex.Lock()
		s.clients[client.GetUUID()] = client
		s.mutex.Unlock()

		if s.onNewClient != nil {
			s.onNewClient(client)
		}

		s.wg.Add(1)
		go s.handleClient(client)
	}
}

func (s *TCPServer) handleClient(client *Client) {
	defer client.conn.Close()
	defer s.wg.Done()

	go func() {
		for message := range client.send {
			client.conn.Write(message)
		}
	}()

	buffer := make([]byte, 4096)
	for {
		n, err := client.conn.Read(buffer)
		if s.IsStopped() {
			if s.onClientClosed != nil {
				s.onClientClosed(client, nil)
			}
			return
		}
		if err != nil {
			close(client.send)
			s.mutex.Lock()
			delete(s.clients, client.GetUUID())
			s.mutex.Unlock()

			if s.onClientClosed != nil {
				s.onClientClosed(client, err)
			}
			return
		}

		if s.onNewMessage != nil {
			s.onNewMessage(client, buffer[:n])
		}
	}
}
