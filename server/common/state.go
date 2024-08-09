package common

type State uint8

const (
	Handshaking State = 0
	Status      State = 1
	Login       State = 2
	Transfer    State = 3
)
