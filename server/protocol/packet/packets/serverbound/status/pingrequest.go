package status

type PingRequestPacket struct {
	Payload uint64 `mc:"long"`
}

func (PingRequestPacket) ID() uint32 {
	return 0x01
}
