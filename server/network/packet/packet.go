package packet

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Packet interface {
	ID() uint32
}

type FieldInfo struct {
	Type   string
	MaxLen int
}

func ParseMCTag(field reflect.StructField) FieldInfo {
	tag := field.Tag.Get("mc")
	parts := strings.Split(tag, ",")

	info := FieldInfo{
		Type:   parts[0],
		MaxLen: 32767,
	}

	if len(parts) > 1 {
		for _, part := range parts[1:] {
			if strings.HasPrefix(part, "max=") {
				maxLen, _ := strconv.Atoi(strings.TrimPrefix(part, "max="))
				info.MaxLen = maxLen
			}
		}
	}

	return info
}

func EncodePacket(p Packet) ([]byte, error) {
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var buffer []byte

	// Encode packet ID
	id := p.ID()
	idBytes := EncodeVarInt(id)
	buffer = append(buffer, idBytes...)

	// Encode fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		info := ParseMCTag(fieldType)

		var fieldBytes []byte
		var err error

		switch info.Type {
		case "varint":
			fieldBytes = EncodeVarInt(uint32(field.Uint()))
		case "string":
			str := field.String()
			if info.MaxLen > 0 && len(str) > info.MaxLen {
				return nil, errors.New("string too long")
			}
			fieldBytes = append(EncodeVarInt(uint32(len(str))), []byte(str)...)
		case "ushort":
			fieldBytes = make([]byte, 2)
			binary.BigEndian.PutUint16(fieldBytes, uint16(field.Uint()))
		case "long":
			fieldBytes = make([]byte, 8)
			binary.BigEndian.PutUint64(fieldBytes, field.Uint())
		case "JSON":
			jsonBytes, err := json.Marshal(field.Interface())
			if err != nil {
				return nil, err
			}
			fieldBytes = append(EncodeVarInt(uint32(len(jsonBytes))), jsonBytes...)
		default:
			return nil, errors.New("unsupported field type")
		}

		buffer = append(buffer, fieldBytes...)
	}

	return buffer, nil
}

func DecodePacket(data []byte, p Packet) error {
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Skip packet ID
	_, n := DecodeVarInt(data)
	data = data[n:]

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		info := ParseMCTag(fieldType)

		var value interface{}
		var err error

		switch info.Type {
		case "varint":
			value, n = DecodeVarInt(data)
			data = data[n:]
		case "string":
			length, n := DecodeVarInt(data)
			data = data[n:]
			if info.MaxLen > 0 && int(length) > info.MaxLen {
				return errors.New("string too long")
			}
			value = string(data[:length])
			data = data[length:]
		case "ushort":
			value = binary.BigEndian.Uint16(data[:2])
			data = data[2:]
		case "long":
			value = binary.BigEndian.Uint64(data[:8])
			data = data[8:]
		case "JSON":
			length, n := DecodeVarInt(data)
			data = data[n:]
			err = json.Unmarshal(data[:length], field.Addr().Interface())
			if err != nil {
				return err
			}
			data = data[length:]
			continue // Skip setting the value, as it's already unmarshaled
		default:
			return errors.New("unsupported field type")
		}

		field.Set(reflect.ValueOf(value))
	}

	return nil
}
