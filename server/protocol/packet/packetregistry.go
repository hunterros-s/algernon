package packet

import (
	"log"

	"github.com/hunterros-s/algernon/server/common"
	"github.com/hunterros-s/algernon/server/protocol/io"
)

type decoder func(*io.Reader) (common.ServerboundPacket, error)

var serverbound_packet_index = map[common.State]map[uint32]decoder{}

// register adds a decoder function for a specific packet ID and state.
func RegisterDecoder(state common.State, packetID uint32, dec decoder) {
	log.Printf("Registering packet with id: %d, state: %d\n", state, packetID)
	if _, exists := serverbound_packet_index[state]; !exists {
		serverbound_packet_index[state] = make(map[uint32]decoder)
	}
	serverbound_packet_index[state][packetID] = dec
}

func GetDecoder(state common.State, packetID uint32) (decoder, bool) {
	stateMap, stateExists := serverbound_packet_index[state]
	if !stateExists {
		return nil, false
	}

	dec, exists := stateMap[packetID]
	return dec, exists
}
