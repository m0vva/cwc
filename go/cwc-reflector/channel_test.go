package main

import (
	"testing"
	"gotest.tools/assert"
)

func TestNewChannel(t *testing.T) {
	channel := NewChannel(33)
	assert.DeepEqual(t,channel.ChannelId, uint16(33))
	assert.Equal(t, len(channel.Subscribers), 0)
	assert.Equal(t, len(channel.Addresses), 0)
}

func TestGetChannel(t *testing.T) {
	channel1 := GetChannel(21)
	channel2 := GetChannel(21)
	assert.DeepEqual(t, channel1, channel2)
}
