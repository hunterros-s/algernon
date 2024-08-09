package supervisor

import (
	"github.com/hunterros-s/algernon/server/common"
	"github.com/hunterros-s/algernon/server/protocol/packet/packets/serverbound/handshaking"
	"github.com/rs/zerolog"
)

type Supervisor struct {
	incoming chan common.IncomingEntry
	logger   zerolog.Logger

	// players  map[uuid.UUID]*Player
	// world    *World
}

func NewSupervisor(logger zerolog.Logger) *Supervisor {
	return &Supervisor{
		incoming: make(chan common.IncomingEntry),
		logger:   logger,
	}
}

// Handle adds a new packet entry to the channel.
func (sv *Supervisor) Handle(entry common.IncomingEntry) {
	sv.incoming <- entry
}

// Start initializes the goroutine that handles incoming entries.
func (sv *Supervisor) Start() {
	go sv.supervise()
}

func (sv *Supervisor) Stop() {
	close(sv.incoming)
}

func (sv *Supervisor) supervise() {
	for {
		entry, ok := <-sv.incoming
		if !ok {
			return
		}
		switch packet := entry.Packet.(type) {
		case *handshaking.HandshakePacket:
			sv.logger.Info().Msg("Handshake packet recieved")
		default:
			sv.logger.Warn().Int("packet id", int(packet.MCPacketID())).Msg("Unknown packet type")
		}
		sv.logger.Info().Str("client uuid", entry.Client.GetUUID().String()).Send()
		sv.logger.Info().Int("packet id", int(entry.Packet.MCPacketID())).Send()
		sv.logger.Info().Str("packet uid", entry.Packet.PacketUID()).Send()
	}
}
