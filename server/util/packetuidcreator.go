package util

import (
	"fmt"

	"github.com/google/uuid"
)

// var uid = uuid.New()

// func (HandshakePacket) PacketUID() []byte {
// 	return uid[:]
// }

func GetPacketUID(p interface{}) string {
	name := fmt.Sprintf("%T", p)
	uuid := uuid.New().String()
	return fmt.Sprintf("%s-%s", name, uuid)
}
