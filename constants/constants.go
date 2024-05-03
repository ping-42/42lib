package constants

import "time"

const (
	ProtoTCP = 1
	ProtoUDP = 2
	//
	TelemetryMonitorPeriod          = time.Minute * 5  // the period in which the sensor will send telemtry data to the server
	TelemetryMonitorPeriodThreshold = time.Second * 40 // the ttl for active sensors will be TelemetryMonitorPeriod+TelemetryMonitorPeriodThreshold - afrer this the sensor will be considered as offline
	//
)
