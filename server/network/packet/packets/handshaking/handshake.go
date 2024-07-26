package handshake

type HandshakePacket struct {
	ProtocolVersion uint32 `mc:"varint"`
	ServerAddress   string `mc:"string,max=255"`
	ServerPort      uint32 `mc:"ushort"`
	NextState       uint32 `mc:"varint"`
}

func (HandshakePacket) ID() uint32 {
	return 0x00
}
