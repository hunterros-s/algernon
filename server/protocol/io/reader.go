package io

import (
	"bytes"

	"github.com/google/uuid"
	"github.com/hunterros-s/algernon/text"
)

type Reader struct {
	buffer *bytes.Buffer
	err    error
}

func NewReader(buf []byte) *Reader {
	return &Reader{buffer: bytes.NewBuffer(buf)}
}

func (r *Reader) Err() error {
	return r.err
}

func (r *Reader) Bytes() []byte {
	return r.buffer.Bytes()
}

func (r *Reader) ReadBool() bool {
	if r.err != nil {
		return false
	}
	var value bool
	value, r.err = ReadBool(r.buffer)
	return value
}

func (r *Reader) ReadByteInt8() int8 {
	if r.err != nil {
		return 0
	}
	var value int8
	value, r.err = ReadByte(r.buffer)
	return value
}

func (r *Reader) ReadUbyte() byte {
	if r.err != nil {
		return 0
	}
	var value byte
	value, r.err = ReadUbyte(r.buffer)
	return value
}

func (r *Reader) ReadShort() int16 {
	if r.err != nil {
		return 0
	}
	var value int16
	value, r.err = ReadShort(r.buffer)
	return value
}

func (r *Reader) ReadUshort() uint16 {
	if r.err != nil {
		return 0
	}
	var value uint16
	value, r.err = ReadUshort(r.buffer)
	return value
}

func (r *Reader) ReadInt() int32 {
	if r.err != nil {
		return 0
	}
	var value int32
	value, r.err = ReadInt(r.buffer)
	return value
}

func (r *Reader) ReadLong() int64 {
	if r.err != nil {
		return 0
	}
	var value int64
	value, r.err = ReadLong(r.buffer)
	return value
}

func (r *Reader) ReadFloat() float32 {
	if r.err != nil {
		return 0
	}
	var value float32
	value, r.err = ReadFloat(r.buffer)
	return value
}

func (r *Reader) ReadDouble() float64 {
	if r.err != nil {
		return 0
	}
	var value float64
	value, r.err = ReadDouble(r.buffer)
	return value
}

func (r *Reader) ReadString() string {
	if r.err != nil {
		return ""
	}
	var value string
	value, r.err = ReadString(r.buffer)
	return value
}

func (r *Reader) ReadJSONTextComponent() text.TextComponent {
	if r.err != nil {
		return text.TextComponent{}
	}
	var value text.TextComponent
	value, r.err = ReadJSONTextComponent(r.buffer)
	return value
}

func (r *Reader) ReadIdentifier() string {
	if r.err != nil {
		return ""
	}
	var value string
	value, r.err = ReadIdentifier(r.buffer)
	return value
}

func (r *Reader) ReadVarInt() int32 {
	if r.err != nil {
		return 0
	}
	var value int32
	value, r.err = ReadVarInt(r.buffer)
	return value
}

func (r *Reader) ReadVarLong() int64 {
	if r.err != nil {
		return 0
	}
	var value int64
	value, r.err = ReadVarLong(r.buffer)
	return value
}

func (r *Reader) ReadPosition() (int32, int32, int32) {
	if r.err != nil {
		return 0, 0, 0
	}
	var x, y, z int32
	x, y, z, r.err = ReadPosition(r.buffer)
	return x, y, z
}

func (r *Reader) ReadUUID() uuid.UUID {
	if r.err != nil {
		return uuid.UUID{}
	}
	var value uuid.UUID
	value, r.err = ReadUUID(r.buffer)
	return value
}

func (r *Reader) ReadBitSet() []int64 {
	if r.err != nil {
		return nil
	}
	var value []int64
	value, r.err = ReadBitSet(r.buffer)
	return value
}

func (r *Reader) ReadFixedBitSet(bits int32) []byte {
	if r.err != nil {
		return nil
	}
	var value []byte
	value, r.err = ReadFixedBitSet(r.buffer, bits)
	return value
}

func (r *Reader) ReadByteArray() []byte {
	if r.err != nil {
		return nil
	}
	var value []byte
	value, r.err = ReadByteArray(r.buffer)
	return value
}

func (r *Reader) ReadFixedByteArray(length int) []byte {
	if r.err != nil {
		return nil
	}
	var value []byte
	value, r.err = ReadFixedByteArray(r.buffer, length)
	return value
}
