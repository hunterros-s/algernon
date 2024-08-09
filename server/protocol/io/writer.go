package io

import (
	"github.com/google/uuid"
	"github.com/hunterros-s/algernon/text"
)

type Writer struct {
	buffer []byte
	err    error
}

func NewWriter() *Writer {
	return &Writer{
		buffer: make([]byte, 0),
		err:    nil,
	}
}

func (w *Writer) Err() error {
	return w.err
}

func (w *Writer) Bytes() []byte {
	return w.buffer
}

func (w *Writer) WriteBool(b bool) *Writer {
	if w.err != nil {
		return nil
	}
	buf, err := WriteBool(b)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteByteInt8(i int8) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteByte(i)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteUbyte(i uint8) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteUbyte(i)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteShort(i int16) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteShort(i)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteUshort(i uint16) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteUshort(i)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteInt(i int32) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteInt(i)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteLong(i int64) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteLong(i)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteFloat(f float32) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteFloat(f)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteDouble(f float64) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteDouble(f)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteString(s string) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteString(s)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteJSONTextComponent(comp text.TextComponent) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteJSONTextComponent(comp)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteIdentifier(s string) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteIdentifier(s)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteVarInt(value int32) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteVarInt(value)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteVarLong(value int64) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteVarLong(value)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WritePosition(x, y, z int32) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WritePosition(x, y, z)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteUUID(u uuid.UUID) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteUUID(u)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteBitSet(data []int64) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteBitSet(data)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteFixedBitSet(data []byte) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteFixedBitSet(data)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteByteArray(s []byte) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteByteArray(s)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}

func (w *Writer) WriteFixedByteArray(s []byte) *Writer {
	if w.err != nil {
		return w
	}
	buf, err := WriteFixedByteArray(s)
	if err != nil {
		w.err = err
		return w
	}
	w.buffer = append(w.buffer, buf...)
	return w
}
