package status

type StatusRequestPacket struct {
	// no fields
}

func (StatusRequestPacket) ID() uint32 {
	return 0x00
}
