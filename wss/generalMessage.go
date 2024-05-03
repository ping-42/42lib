package wss

type GeneralMessage struct {
	MessageGeneralType MessageGeneralType
}

// WssMessageType each message sent via wss must have WssMessageType
type MessageGeneralType int

const (
	// possible WSS messages types
	MessageTypeTaskResult MessageGeneralType = iota
	MessageTypeTelemtry
)

func (s MessageGeneralType) String() string {
	switch s {
	case MessageTypeTaskResult:
		return "MessageTypeTaskResult"
	case MessageTypeTelemtry:
		return "MessageTypeTelemtry"
	}
	return "unknown"
}
