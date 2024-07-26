package packet

import (
	"context"
	"encoding/hex"

	"github.com/hunterros-s/algernon/logger"
	"github.com/hunterros-s/algernon/server/network/packet/codec"
	"github.com/hunterros-s/algernon/server/network/session"
)

type PacketHandler struct{}

func (p_h *PacketHandler) ReadPacket(ctx context.Context, session *session.Session, byteslice []byte, log logger.Logger) {
	select {
	case <-ctx.Done():
		// The context has been cancelled, so we should stop processing
		log.Info().Msg("Context cancelled, stopping packet processing")
		return
	default:
		// Continue processing the packet
	}

	reader := codec.NewReader(byteslice)

	packet_length := reader.ReadVarInt()
	packet_id := reader.ReadVarInt()

	if reader.Error() != nil {
		log.Error().Err(reader.Error()).Msg("Failed to read packet")
		return
	}

	log.Debug().Int("length", int(packet_length)).Int("packet id", int(packet_id)).Str("data", hex.EncodeToString(reader.Bytes())).Msg("Received packet")

	// You could add more context-aware processing here
	// For example, you could pass the context to other functions that might need it
	err := p_h.processPacket(ctx, session, packet_id, reader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process packet")
	}
}

func (p_h *PacketHandler) processPacket(ctx context.Context, session *session.Session, packetID int32, reader *codec.Reader) error {
	// This is where you would add your packet processing logic
	// You can use the context here for timeouts or cancellation
	return nil
}

func NewPacketHandler() *PacketHandler {
	return &PacketHandler{}
}
