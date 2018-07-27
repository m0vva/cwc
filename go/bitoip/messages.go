package bitoip

// conservative UDP payload in bytes
const MaxMessageSizeInBytes = 400

type (
	MessageVerb = uint8
	ChannelId = uint16
	CarrierKey = uint16
	Callsign = string
)

const (
	EnumerateChannels MessageVerb = 0x90
	ListChannels MessageVerb = 0x91
	TimeSync MessageVerb = 0x92
	TimeSyncResponse MessageVerb = 0x93
	ListenRequest MessageVerb = 0x94
	ListenConfirm MessageVerb = 0x95
	Unlisten MessageVerb = 0x96
	KeyValue MessageVerb = 0x81
	CarrierEvent MessageVerb = 0x82
)

// List of Channels
const MaxChannelsPerMessage int = (MaxMessageSizeInBytes - 1) / 2
type ListChannelsPayload struct {
	channels [MaxChannelsPerMessage]uint16
}

// TimeSync
type TimeSyncPayload struct {
	currentTime uint64
}

type TimeSyncResponsePayload struct {
	givenTime uint64
	currentTime uint64
}

type ListenRequestPayload struct {
	channel ChannelId
	callsign string
}

type ListenConfirmPayload struct {
	channel ChannelId
	carrierKey CarrierKey
}

type UnlistenPayload struct {
	channel ChannelId
	carrierKey CarrierKey
}

type KeyValuePayload struct {
	channel ChannelId
	carrierKey CarrierKey
	key string
	value string
}

type BitEvent uint8

const (
	BitOn BitEvent = 0x01
	BitOff BitEvent = 0x00
)

// slightly random
const MaxBitEvents = (MaxMessageSizeInBytes - 10) / 10

type CarrierBitEvent struct {
	timestamp uint64
	bitEvent BitEvent
}

type CarrierEventPayload struct {
	channel ChannelId
	carrierKey [MaxBitEvents]CarrierKey
}