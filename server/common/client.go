package common

import "github.com/google/uuid"

type Client interface {
	Send(b []byte)
	GetState() State
	GetUUID() uuid.UUID
}
