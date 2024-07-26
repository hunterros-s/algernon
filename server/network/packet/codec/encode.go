package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/hunterros-s/algernon/text"
)

// https://wiki.vg/Protocol#Type:Boolean
func WriteBool(b bool) ([]byte, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// https://wiki.vg/Protocol#Type:Byte
func WriteByte(i int8) ([]byte, error) {
	return []byte{byte(i)}, nil
}

// https://wiki.vg/Protocol#Type:Unsigned_Byte
func WriteUbyte(i uint8) ([]byte, error) {
	return []byte{i}, nil
}

// https://wiki.vg/Protocol#Type:Short
func WriteShort(i int16) ([]byte, error) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(i))
	return buf, nil
}

// https://wiki.vg/Protocol#Type:Unsigned_Short
func WriteUshort(i uint16) ([]byte, error) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, i)
	return buf, nil
}

// https://wiki.vg/Protocol#Type:Int
func WriteInt(i int32) ([]byte, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf, nil
}

// https://wiki.vg/Protocol#Type:Long
func WriteLong(i int64) ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf, nil
}

// https://wiki.vg/Protocol#Type:Float
func WriteFloat(f float32) ([]byte, error) {
	return WriteInt(int32(math.Float32bits(f)))
}

// https://wiki.vg/Protocol#Type:Double
func WriteDouble(f float64) ([]byte, error) {
	return WriteLong(int64(math.Float64bits(f)))
}

// https://wiki.vg/Protocol#Type:String
func WriteString(s string) ([]byte, error) {
	lenBytes, err := WriteVarInt(int32(len(s)))
	if err != nil {
		return nil, err
	}
	return append(lenBytes, []byte(s)...), nil
}

// https://wiki.vg/Protocol#Type:JSON_Text_Component
func WriteJSONTextComponent(comp text.TextComponent) ([]byte, error) {
	d, err := json.Marshal(comp)
	if err != nil {
		return nil, err
	}
	return WriteByteArray(d)
}

// https://wiki.vg/Protocol#Type:Identifier
func WriteIdentifier(s string) ([]byte, error) {
	if len(s) > 32767 {
		return nil, fmt.Errorf("expected identifier len to be <= 32767, got %d", len(s))
	}
	return WriteString(s)
}

// https://wiki.vg/Protocol#Type:VarInt
func WriteVarInt(value int32) ([]byte, error) {
	var buf bytes.Buffer
	ux := uint32(value)
	for ux >= 0x80 {
		if err := buf.WriteByte(byte(ux&0x7F) | 0x80); err != nil {
			return nil, err
		}
		ux >>= 7
	}
	if err := buf.WriteByte(byte(ux)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// https://wiki.vg/Protocol#Type:VarLong
func WriteVarLong(value int64) ([]byte, error) {
	var buf bytes.Buffer
	ux := uint64(value)
	for ux >= 0x80 {
		if err := buf.WriteByte(byte(ux&0x7F) | 0x80); err != nil {
			return nil, err
		}
		ux >>= 7
	}
	if err := buf.WriteByte(byte(ux)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// https://wiki.vg/Protocol#Type:Position
func WritePosition(x, y, z int32) ([]byte, error) {
	value := ((int64(x) & 0x3FFFFFF) << 38) | ((int64(z) & 0x3FFFFFF) << 12) | (int64(y) & 0xFFF)
	return WriteLong(value)
}

// https://wiki.vg/Protocol#Type:UUID
func WriteUUID(u uuid.UUID) ([]byte, error) {
	return u[:], nil
}

// https://wiki.vg/Protocol#Type:BitSet
func WriteBitSet(data []int64) ([]byte, error) {
	lenBytes, err := WriteVarInt(int32(len(data)))
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(lenBytes)
	for _, l := range data {
		longBytes, err := WriteLong(l)
		if err != nil {
			return nil, err
		}
		buf.Write(longBytes)
	}
	return buf.Bytes(), nil
}

// https://wiki.vg/Protocol#Type:Fixed_BitSet
func WriteFixedBitSet(data []byte) ([]byte, error) {
	return data, nil
}

// https://wiki.vg/Protocol#Type:Byte_Array
func WriteByteArray(s []byte) ([]byte, error) {
	lenBytes, err := WriteVarInt(int32(len(s)))
	if err != nil {
		return nil, err
	}
	return append(lenBytes, s...), nil
}

func WriteFixedByteArray(s []byte) ([]byte, error) {
	return s, nil
}
