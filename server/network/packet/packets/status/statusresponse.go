package status

type StatusResponsePacket struct {
	JSONResponse string `mc:"JSON"`
}

func (StatusResponsePacket) ID() uint32 {
	return 0x00
}
