package bitoip

import (
	"testing"
		"gotest.tools/assert"
		)

func TestTimeSyncMessageEncode(t *testing.T) {
	verb := TimeSync
	payload := TimeSyncPayload{1}

	myBytes := EncodePayload(verb, payload)
	assert.DeepEqual(t, myBytes, []uint8{146, 0, 0, 0, 0, 0, 0, 0, 1})
}

func TestEnumerateChannels(t *testing.T) {
	verb := EnumerateChannels
	payload := Payload(nil)

	myBytes := EncodePayload(verb, payload)
	assert.DeepEqual(t, myBytes, []uint8{EnumerateChannels})
}

func TestEnMessageEncode(t *testing.T) {
	verb := ListChannels
	payload := ListChannelsPayload{[MaxChannelsPerMessage]uint16{1,3,2}}

	myBytes := EncodePayload(verb, payload)
	assert.DeepEqual(t, myBytes[0:7], []uint8{ListChannels, 0, 1, 0, 3, 0, 2})
}
