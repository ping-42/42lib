package wss

import (
	"net"

	"github.com/google/uuid"
)

// SensorConnection define ws client connection
type SensorConnection struct {
	// Uuid unique id per each connection
	ConnectionId uuid.UUID
	// Connection ws connection
	// we do not need this field storing it to Redis
	Connection    net.Conn `json:"-"`
	SensorId      uuid.UUID
	SensorVersion string
}
