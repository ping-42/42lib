package models

import (
	"time"

	"github.com/google/uuid"
)

type TsHostNetworkStat struct {
	Time     time.Time `gorm:"type:TIMESTAMPTZ;"`
	SensorID uuid.UUID `gorm:"type:uuid;"`
}

type TsNetworkInterfaceStat struct {
	NetworkStatID uuid.UUID `gorm:"type:uuid;"`
	InterfaceName string
	BytesSent     uint64
	BytesRecv     uint64
	PacketsSent   uint64
	PacketsRecv   uint64
}

func (TsNetworkInterfaceStat) TableName() string {
	return "ts_network_interface_stats"
}
