package handshake

import (
	"bytes"
	"fmt"

	"github.com/hunterros-s/algernon/server/network/packet/codec"
)

type HandshakePacket struct {
	ProtocolVersion uint32 `mc:"varint"`
	ServerAddress   string `mc:"string,max=255"`
	ServerPort      uint32 `mc:"ushort"`
	NextState       uint32 `mc:"varint"`
}

func (HandshakePacket) ID() uint32 {
	return 0x00
}

func (pkt *HandshakePacket) Encode() (*bytes.Buffer, error) {
	var buffer bytes.Buffer

	protocol_version, err := codec.WriteVarInt(int32(pkt.ProtocolVersion))
	if err != nil {
		return nil, fmt.Errorf("error writing protocol version to varint:%v", err.Error())
	}
	if _, err = buffer.Write(protocol_version); err != nil {
		return nil, fmt.Errorf("error writing protocol version to buffer:%v", err.Error())
	}

	server_address, err := codec.WriteString(pkt.ServerAddress)
	if err != nil {
		return nil, fmt.Errorf("error writing server address to string:%v", err.Error())
	}
	if _, err = buffer.Write(server_address); err != nil {
		return nil, fmt.Errorf("error writing server address to buffer:%v", err.Error())
	}

	server_port, err := codec.WriteUshort(uint16(pkt.ServerPort))
	if err != nil {
		return nil, fmt.Errorf("error writing server port to varint:%v", err.Error())
	}
	if _, err = buffer.Write(server_port); err != nil {
		return nil, fmt.Errorf("error writing server port to buffer:%v", err.Error())
	}

	next_state, err := codec.WriteVarInt(int32(pkt.NextState))
	if err != nil {
		return nil, fmt.Errorf("error writing next state to varint:%v", err.Error())
	}
	if _, err = buffer.Write(next_state); err != nil {
		return nil, fmt.Errorf("error writing next state to buffer:%v", err.Error())
	}

	return &buffer, nil
}
