package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func DecodeBool(buf *bytes.Buffer) (bool, error) {
	if buf.Len() == 0 {
		return false, errors.New("buffer is empty")
	}

	b := buf.Next(1)[0]
	return b != 0, nil
}

func DecodeByte(buf *bytes.Buffer) (int8, error) {
	if buf.Len() == 0 {
		return 0, errors.New("buffer is empty")
	}

	b := buf.Next(1)[0]
	return int8(b), nil
}

func DecodeUbyte(buf *bytes.Buffer) (byte, error) {
	if buf.Len() == 0 {
		return 0, errors.New("buffer is empty")
	}

	b := buf.Next(1)[0]
	return b, nil
}

func DecodeShort(buf *bytes.Buffer) (int16, error) {
	if buf.Len() < 2 {
		return 0, errors.New("buffer has insufficient data")
	}

	d := buf.Next(2)
	return int16(d[0])<<8 | int16(d[1]), nil
}

func DecodeUshort(buf *bytes.Buffer) (uint16, error) {
	if buf.Len() < 2 {
		return 0, errors.New("buffer has insufficient data")
	}

	d := buf.Next(2)
	return uint16(d[0])<<8 | uint16(d[1]), nil
}

func DecodeInt(buf *bytes.Buffer) (int32, error) {
	if buf.Len() < 4 {
		return 0, errors.New("buffer has insufficient data")
	}

	d := buf.Next(4)
	return int32(d[0])<<24 | int32(d[1])<<16 | int32(d[2])<<8 | int32(d[3]), nil
}

func DecodeLong(buf *bytes.Buffer) (int64, error) {
	if buf.Len() < 8 {
		return 0, errors.New("buffer has insufficient data")
	}

	d := buf.Next(8)
	return int64(d[0])<<56 | int64(d[1])<<48 | int64(d[2])<<40 | int64(d[3])<<32 |
		int64(d[4])<<24 | int64(d[5])<<16 | int64(d[6])<<8 | int64(d[7]), nil
}

func DecodeFloat(buf *bytes.Buffer) (float32, error) {
	if buf.Len() < 4 {
		return 0, errors.New("buffer has insufficient data")
	}

	d := buf.Next(4)
	bits := uint32(d[0])<<24 | uint32(d[1])<<16 | uint32(d[2])<<8 | uint32(d[3])
	return math.Float32frombits(bits), nil
}

func DecodeDouble(buf *bytes.Buffer) (float64, error) {
	if buf.Len() < 8 {
		return 0, errors.New("buffer has insufficient data")
	}

	d := buf.Next(8)
	bits := uint64(d[0])<<56 | uint64(d[1])<<48 | uint64(d[2])<<40 | uint64(d[3])<<32 |
		uint64(d[4])<<24 | uint64(d[5])<<16 | uint64(d[6])<<8 | uint64(d[7])
	return math.Float64frombits(bits), nil
}

func DecodeString(buf *bytes.Buffer) (string, error) {
	length, err := DecodeVarInt(buf)
	if err != nil {
		return "", err
	}

	if length < 0 || length > 32767*3+3 {
		return "", errors.New("invalid string length")
	}

	if buf.Len() < int(length) {
		return "", errors.New("buffer has insufficient data")
	}

	bytes := buf.Next(int(length))
	return string(bytes), nil
}

func DecodeIdentifier(buf *bytes.Buffer) (string, error) {
	identifier, err := DecodeString(buf)
	if err != nil {
		return "", err
	}

	// Validate the identifier
	if !isValidIdentifier(identifier) {
		return "", fmt.Errorf("invalid identifier format: %s", identifier)
	}

	return identifier, nil
}

func isValidIdentifier(identifier string) bool {
	parts := strings.SplitN(identifier, ":", 2)

	if len(parts) == 1 {
		// If no namespace is provided, only validate the value
		return isValidValue(parts[0])
	} else if len(parts) == 2 {
		// Validate both namespace and value
		return isValidNamespace(parts[0]) && isValidValue(parts[1])
	}

	return false
}

func isValidNamespace(namespace string) bool {
	match, _ := regexp.MatchString("^[a-z0-9.-_]+$", namespace)
	return match
}

func isValidValue(value string) bool {
	match, _ := regexp.MatchString("^[a-z0-9.-_/]+$", value)
	return match
}

func DecodeVarInt(buf *bytes.Buffer) (int32, error) {
	var value int32
	var position int32
	var currentByte byte
	const CONTINUE_BIT byte = 128
	const SEGMENT_BITS byte = 127

	for {
		if buf.Len() == 0 {
			return 0, errors.New("buffer has insufficient data")
		}

		currentByte = buf.Next(1)[0]
		value |= int32((currentByte & SEGMENT_BITS)) << position

		if (currentByte & CONTINUE_BIT) == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return 0, errors.New("VarInt is too big")
		}
	}

	return value, nil
}

func DecodeVarLong(buf *bytes.Buffer) (int64, error) {
	var value int64
	var position int64
	var currentByte byte
	const CONTINUE_BIT byte = 128
	const SEGMENT_BITS byte = 127

	for {
		if buf.Len() == 0 {
			return 0, errors.New("buffer has insufficient data")
		}

		currentByte = buf.Next(1)[0]
		value |= int64((currentByte & SEGMENT_BITS)) << position

		if (currentByte & CONTINUE_BIT) == 0 {
			break
		}

		position += 7

		if position >= 64 {
			return 0, errors.New("VarLong is too big")
		}
	}

	return value, nil
}

func DecodePosition(buf *bytes.Buffer) (x, y, z int32, err error) {
	if buf.Len() < 8 {
		return 0, 0, 0, errors.New("buffer has insufficient data")
	}

	var l int64
	binary.Read(buf, binary.BigEndian, &l)

	x = int32(l >> 38)
	y = int32(l << 52 >> 52)
	z = int32(l << 26 >> 38)

	return x, y, z, nil
}

func DecodeUUID(buf *bytes.Buffer) (uuid.UUID, error) {
	if buf.Len() < 16 {
		return uuid.UUID{}, errors.New("buffer has insufficient data")
	}

	u := make([]byte, 16)
	_, err := buf.Read(u)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.UUID(u), nil
}

func DecodeBitSet(buf *bytes.Buffer) ([]int64, error) {
	length, err := DecodeVarInt(buf)
	if err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, errors.New("invalid BitSet length")
	}

	data := make([]int64, length)

	for i := 0; i < int(length); i++ {
		if buf.Len() < 8 {
			return nil, errors.New("buffer has insufficient data")
		}

		err = binary.Read(buf, binary.BigEndian, &data[i])
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
