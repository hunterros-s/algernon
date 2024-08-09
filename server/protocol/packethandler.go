package protocol

import (
	"fmt"

	"github.com/hunterros-s/algernon/server/common"
	"github.com/hunterros-s/algernon/server/protocol/io"
	"github.com/hunterros-s/algernon/server/protocol/packet"

	// "github.com/hunterros-s/algernon/server/protocol/packet/packets/serverbound"
	"github.com/rs/zerolog"
)

// l.SetOnNewMessage(func(c *listener.Client, p packet.ServerboundPacket) {
// 	cfg.Logger.Info().Int("packet id", int(p.ID())).Send()
// })

// do encoding/decoding shit here
// packet, err := serverbound.ReadUncompressedPacket(client.State, b)

// if err != nil {
// 	l.logger.Error().Err(err).Msg("Packet reading error")
// 	return
// }

// func (c *Client) Send(p packet.ClientboundPacket) error {
// 	data, err := p.Encode()
// 	if err != nil {
// 		return err
// 	}
// 	c.tcpclient.Send(data)
// 	return nil
// }

type PacketHandler func(common.Client, common.ServerboundPacket)

func GetNewMessageCallback(handler PacketHandler, logger zerolog.Logger) func(c common.Client, b []byte) {
	return func(c common.Client, b []byte) {
		packet, err := ReadUncompressedPacket(c.GetState(), b)
		if err != nil {
			logger.Warn().Err(err).Msg("Packet error")
			return
		}

		handler(c, packet)
	}
}

func ReadUncompressedPacket(state common.State, b []byte) (common.ServerboundPacket, error) {
	r := io.NewReader(b)

	_ = r.ReadVarInt() // don't need to use packet length for shit
	packet_id := r.ReadVarInt()

	if r.Err() != nil {
		return nil, r.Err()
	}

	// decoder, ok := serverbound.GetDecoder(state, uint32(packet_id))
	decoder, ok := packet.GetDecoder(state, uint32(packet_id))
	if !ok {
		return nil, fmt.Errorf("unknown packet state: %d, id: %d", state, packet_id)
	}

	packet, err := decoder(r)
	if err != nil {
		return nil, err
	}

	return packet, nil
}
