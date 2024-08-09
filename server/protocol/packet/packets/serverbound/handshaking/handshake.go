package handshaking

import (
	"fmt"

	"github.com/hunterros-s/algernon/server/common"
	"github.com/hunterros-s/algernon/server/protocol/io"
	"github.com/hunterros-s/algernon/server/protocol/packet"
	"github.com/hunterros-s/algernon/server/util"
)

var _ common.ServerboundPacket = (*HandshakePacket)(nil)

type HandshakePacket struct {
	ProtocolVersion int32  `mc:"varint"`
	ServerAddress   string `mc:"string,max=255"`
	ServerPort      uint16 `mc:"ushort"`
	NextState       int32  `mc:"varint"`
}

func (HandshakePacket) MCPacketID() uint32 {
	return 0x00
}

var uid = util.GetPacketUID(HandshakePacket{})

func (HandshakePacket) PacketUID() string {
	return uid
}

func DecodeHandshake(r *io.Reader) (common.ServerboundPacket, error) {
	p_version := r.ReadVarInt()
	s_address := r.ReadString()
	s_port := r.ReadUshort()
	n_state := r.ReadVarInt()

	if r.Err() != nil {
		return nil, fmt.Errorf("error decoding handshake packet: %w", r.Err())
	}

	return &HandshakePacket{
		ProtocolVersion: p_version,
		ServerAddress:   s_address,
		ServerPort:      s_port,
		NextState:       n_state,
	}, nil
}

func init() {
	packet.RegisterDecoder(common.Handshaking, HandshakePacket{}.MCPacketID(), DecodeHandshake)
}
