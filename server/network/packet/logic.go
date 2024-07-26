package packet

import (
	"encoding/hex"
	"net"

	"github.com/hunterros-s/algernon/logger"
)

type PacketHandler struct{}

func (p_h *PacketHandler) ReadPacket(conn net.Conn, byteslice []byte, log logger.Logger) {
	l := len(byteslice)
	log.Info().Int("length", l).Str("msg", hex.EncodeToString(byteslice)).Msgf("Message recieved")
}

func NewPacketHandler() *PacketHandler {
	return &PacketHandler{}
}
