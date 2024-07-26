package status

type PingResponsePacket struct {
	Payload uint64 `mc:"long"`
}

func (PingResponsePacket) ID() uint32 {
	return 0x01
}
