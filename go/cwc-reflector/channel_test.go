package main

import (
	"net"
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

func TestSubscribeWhenNotSubscribed(t *testing.T) {
	channel1 := GetChannel(21)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel1.Subscribe(addr)
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
}

func TestSubscribeWhenSubscribed(t *testing.T) {
	channel1 := GetChannel(21)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel1.Subscribe(addr)
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
	channel1.Subscribe(addr)
	assert.Equal(t, len(channel1.Addresses), 1)
	assert.Equal(t, len(channel1.Subscribers), 1)
}
func TestUnsubscribeWhenSubscribed(t *testing.T) {
	channel2 := GetChannel(22)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel2.Subscribe(addr)
	assert.Equal(t, len(channel2.Addresses), 1)
	assert.Equal(t, len(channel2.Subscribers), 1)
	channel2.Unsubscribe(addr)
	assert.Equal(t, len(channel2.Subscribers), 0)
	assert.Equal(t, len(channel2.Addresses), 0)
}

func TestUnsubscribeWhenNotSubscribed(t *testing.T) {
	channel2 := GetChannel(22)
	addr, _ := net.ResolveUDPAddr("udp", "localhost:2020")
	channel2.Unsubscribe(addr)
	assert.Equal(t, len(channel2.Subscribers), 0)
	assert.Equal(t, len(channel2.Addresses), 0)
}

func TestChannelIds(t *testing.T) {
	GetChannel(21)
	GetChannel(22)
	GetChannel(33)
	assert.DeepEqual(t, ChannelIds(), []uint16{21, 22, 33})
}