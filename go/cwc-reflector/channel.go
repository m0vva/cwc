/*
Copyright (C) 2019 Graeme Sutherland, Nodestone Limited


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"github.com/golang/glog"
	"net"
)
import (
	"../bitoip"
	"time"
)

const LastReservedCarrierKey = 99

type Subscriber struct {
	Key        bitoip.CarrierKeyType
	Address    net.UDPAddr
	LastTx     time.Time
	LastListen time.Time
	Callsign   string
}

type Channel struct {
	ChannelId   bitoip.ChannelIdType
	Subscribers map[bitoip.CarrierKeyType]Subscriber
	Addresses   map[string]Subscriber
	Callsigns   map[string]Subscriber
	LastKey     bitoip.CarrierKeyType
}

type ChannelMap map[uint16]*Channel

var channels = make(ChannelMap)

// Create a new channel
func NewChannel(channelId bitoip.ChannelIdType) Channel {
	return Channel{
		channelId,
		make(map[bitoip.CarrierKeyType]Subscriber),
		make(map[string]Subscriber),
		make(map[string]Subscriber),
		LastReservedCarrierKey,
	}
}

// Return array of channel Ids of existing channels
func ChannelIds() []uint16 {
	keys := make([]uint16, 0, len(channels))
	for k := range channels {
		keys = append(keys, k)
	}
	return keys
}

// Get a channel by channel_id
func GetChannel(channel_id bitoip.ChannelIdType) *Channel {
	if channel, ok := channels[channel_id]; ok {
		return channel
	} else {
		nc := NewChannel(channel_id)
		channels[channel_id] = &nc
		return &nc
	}
}

// Subscribe to this channel
// if already susscribed, then update details and LastTx
func (c *Channel) Subscribe(address net.UDPAddr, callsign string) bitoip.CarrierKeyType {
	glog.Infof("subscribe from: %v", address)
	glog.Infof("channels: %v", channels)
	if subscriber, ok := c.Addresses[address.String()]; ok {
		subscriber.LastListen = time.Now()
		c.Addresses[address.String()] = subscriber
		c.Subscribers[subscriber.Key] = subscriber
		c.Callsigns[callsign] = subscriber
		glog.V(2).Infof("subscribe existing key %d", subscriber.Key)
		return subscriber.Key
	} else {
		c.LastKey += 1
		subscriber := Subscriber{c.LastKey, address, *new(time.Time), time.Now(), callsign}
		c.Subscribers[c.LastKey] = subscriber
		c.Addresses[address.String()] = subscriber
		c.Callsigns[callsign] = subscriber
		glog.V(1).Infof("suscribe new key %d", subscriber.Key)
		return subscriber.Key
	}
}

// Unsubscribe from channel
func (c *Channel) Unsubscribe(address net.UDPAddr) {
	if subscriber, ok := c.Addresses[address.String()]; ok {
		delete(c.Subscribers, subscriber.Key)
		delete(c.Addresses, subscriber.Address.String())
		delete(c.Callsigns, subscriber.Callsign)
	}
}

// Broadcast this carrier event to all on this channel
// and always return to sender (who can ignore if they wish, or can use as net sidetone
func (c *Channel) Broadcast(event bitoip.CarrierEventPayload) {
	txr := c.Subscribers[event.CarrierKey]
	txr.LastTx = time.Now()
	for _, v := range c.Subscribers {
		glog.V(2).Infof("sending to subs %v: %v", v.Address, event)
		bitoip.UDPTx(bitoip.CarrierEvent, event, &v.Address)
	}

}

// Check through for subscribers that we haven't seem for a while
// and remove them.
func SuperviseChannels(t time.Time, timeout time.Duration) int {
	removed := 0
	for _, channel := range channels {
		for key, sub := range channel.Subscribers {
			if t.Sub(sub.LastListen) > timeout {
				delete(channel.Subscribers, key)
				removed += 1
				for add, sub := range channel.Addresses {
					if sub.Key == key {
						delete(channel.Addresses, add)
					}
				}
				for call, sub := range channel.Callsigns {
					if sub.Key == key {
						delete(channel.Callsigns, call)
					}
				}
			}
		}
	}

	return removed
}
