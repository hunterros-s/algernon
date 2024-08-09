package listener

import (
	"github.com/google/uuid"
	"github.com/hunterros-s/algernon/server/common"
	"github.com/hunterros-s/algernon/tcpserver"
	"github.com/rs/zerolog"
)

type client struct {
	tcpclient *tcpserver.Client
	Logger    zerolog.Logger
	state     common.State
}

func newClient(c *tcpserver.Client, logger zerolog.Logger) *client {
	return &client{
		tcpclient: c,
		Logger:    logger.With().Str("client address", c.GetIP()).Logger(),
		state:     common.Handshaking,
	}
}

func (c *client) Send(data []byte) {
	c.tcpclient.Send(data)
}

func (c *client) GetState() common.State {
	return c.state
}

func (c *client) GetUUID() uuid.UUID {
	return c.tcpclient.GetUUID()
}
