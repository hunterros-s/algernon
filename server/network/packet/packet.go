package packet

import "bytes"

type Packet interface {
	ID() uint32
	Encode() (*bytes.Buffer, error)
}
