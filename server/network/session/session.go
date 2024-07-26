package session

import "net"

type State uint8

const (
	Handshaking State = 0
	Status      State = 1
	Login       State = 2
	Transfer    State = 3
)

type Session struct {
	connection net.Conn
	state      State
}

func (s *Session) Conn() net.Conn {
	return s.connection
}

func (s *Session) GetState() State {
	return s.state
}

func (s *Session) SetState(i State) {
	s.state = i
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		connection: conn,
		state:      Handshaking,
	}
}
