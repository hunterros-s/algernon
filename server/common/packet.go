package common

type Packet interface {
	MCPacketID() uint32
	PacketUID() string
}

// all serverbound packets need decoders
type ServerboundPacket interface {
	MCPacketID() uint32
	PacketUID() string
}

// all clientbound packets need an encoder to go from packet -> bytes
type ClientboundPacket interface {
	MCPacketID() uint32
	PacketUID() string
	Encode() ([]byte, error)
}
