package bitoip

import (
	"encoding/binary"
	"bytes"
	"log"
	"reflect"
	)

// conservative UDP payload in bytes
const MaxMessageSizeInBytes = 400
const CallsignSize = 16
var byteOrder = binary.BigEndian

type (
	MessageVerb = byte
	ChannelIdType = uint16
	CarrierKeyType = uint16
	Callsign = [CallsignSize]byte
	Payload = interface {}
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
	channel ChannelIdType
	callsign [16]byte
}

type ListenConfirmPayload struct {
	channel ChannelIdType
	carrierKey CarrierKeyType
}

type UnlistenPayload struct {
	channel ChannelIdType
	carrierKey CarrierKeyType
}

type KeyValuePayload struct {
	channel ChannelIdType
	carrierKey CarrierKeyType
	key [8]byte
	value [16]byte
}

type BitEvent uint8

const (
	BitOn BitEvent = 0x01
	BitOff BitEvent = 0x00
	LastEvent BitEvent = 0x80 // high bit set to indicate last one
)

// slightly random
const MaxBitEvents = (MaxMessageSizeInBytes - 14) / 5
const MaxNsPerCarrierEvent = 2 ^ 32

// Offset allows for about 4 seconds of offset
type CarrierBitEvent struct {
	TimeOffset uint32
	BitEvent BitEvent
}

type CarrierEventPayload struct {
	Channel ChannelIdType
	CarrierKey CarrierKeyType
	StartTimeStamp uint64
	BitEvents [MaxBitEvents]CarrierBitEvent
}

var messagePayload = map[MessageVerb]reflect.Type {
	EnumerateChannels: nil,
	ListChannels: reflect.TypeOf(ListChannelsPayload{}),
	TimeSync: reflect.TypeOf(TimeSyncPayload{}),
	TimeSyncResponse: reflect.TypeOf(TimeSyncResponsePayload{}),
	ListenRequest: reflect.TypeOf(ListenRequestPayload{}),
	ListenConfirm: reflect.TypeOf(ListenConfirmPayload{}),
	Unlisten: reflect.TypeOf(UnlistenPayload{}),
	KeyValue: reflect.TypeOf(KeyValuePayload{}),
	CarrierEvent: reflect.TypeOf(CarrierEventPayload{}),
}


func EncodePayload(verb MessageVerb, payload Payload) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(verb)
	if payload != nil {
		err := binary.Write(buf, byteOrder, payload)
		if err != nil {
			log.Fatalf("Bad message encode for %T %v", payload, err)
		}
	}
	return buf.Bytes()
}

func DecodePacket(lineBuffer []byte) (MessageVerb, interface{}) {
	verb := MessageVerb(lineBuffer[0])
	var payloadObj interface{} = nil

	switch verb {
	case EnumerateChannels:
		break;
	case ListChannels:
		payloadObj = new(ListChannelsPayload)
	case TimeSync:
		payloadObj = new(TimeSyncPayload)
	case TimeSyncResponse:
		payloadObj = new(TimeSyncResponsePayload)
	case ListenRequest:
		payloadObj = new(ListenRequestPayload)
	case ListenConfirm:
		payloadObj = new(ListenConfirmPayload)
	case Unlisten:
		payloadObj = new(UnlistenPayload)
	case KeyValue:
		payloadObj = new(KeyValuePayload)
	case CarrierEvent:
		payloadObj = new(CarrierEventPayload)

	}
	buffer := bytes.NewReader(lineBuffer[1:])
	err := binary.Read(buffer, byteOrder, payloadObj)

	return verb, payloadObj
}
