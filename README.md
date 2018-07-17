# bitoip

## What

bit-over-ip: Low level protocol for sending timed on-off state over TCP/IP or UDP/IP

## Why

I'm looking for a way to transmit hand-keyed and keyer-keyed morse over the
internet.  Ideally, this is via a toggled line on a serial port or software
generated from morse keying or morse detecting software.

This becomes generated *carrier* that is a bunch of on and off events, 
with strict timing maintained so it can be re-constructed at
the far end.  Packers of on/off events are sent via UDP, as is this whole protocol.  Data
can get lost, but we can deal with that.

Carriers are grouped into a *channel* which has a name, and some control interface.

Channels are hosted on a *server*.

##  Entities in detail

Here is some more detail on what the logical parts are:

A **carrier** is a set of bit on/off events for a single bit.  A carrier has:
 - an id (which is used to differentiate different carrier in a channel)
 - a series of on and off events and control messages
 - some sort of sender id (like a callsign)

Note that a carrier should always end up off, not stuck on.  So timed on/off events must always end with an off.  That means the way we collect and send bit data has to
take that into account.

A **channel** is a named grouping of carriers.  A server may support one or more channels.  A channel has:
 - a control interface to allow:
    - enumeration of carriers
    - subscription and unsubscription
 - a name.  Should be unique in a server
 - combining with a server IP address and port, should be possible to make a channel URI

 Given that a channel can have multiple carriers, the receiving end will need to differentiate them when making them audible by
 using tone offsets or something similar, so indvidual signals can be recognised. Note that a channel originates at a hub.

A **hub** can receive and publish one or more channels by communicating with nodes and (later) possibly other hubs.
There is no hub - to -hub routing mechanism proposed here yet.

A **node** can register interest in channels with a hub and add carriers to a channel.

## How this works

A spoke would connect to a hub and probably get an enumeration of channels.  The client can then
subscribe to a channel to receive carrier packets related to that channel.  The subscription results in time
sync being established for that channel, and then the client will receive packets relating to that channel.  The client
can also transmit packets relating to that channel based on the time sync.

There can be multiple carriers per channel, so it is up to the client to make sense of these.


## Protocol

## TODO
More protocol tidy up

### Channel access

The protocol is composed of packets containing timed on and off information for a bit, with the
necessary time offsets from the start of the stream.

There is a facility to communicate basic key/value ascii pairs.

Generally, top bit set means that is a control value.

multi-byte values are big-endian

## Enumerate channels

Enumerate channels at this hub

EN (Enumerate) == 0x90

## Enumerate list

Response to enumerate channels available:

EL (Enumerate List) == 0xA0

0x91, channel_no (2 bytes), channel_no (2 bytes), ... , 0x00

## Server time sync

There are a pair of time sync events to allow nodes and hubs to calculate lateny in the circuit and adjust for
it if needed.

### Initiate time sync
0x92 current_time [timestamp size]

### Respond to time sync
0x93 given_time_or 0 (from prev 0x91) [8 byte, unix nanoseconds], current_time [8byte unix nanoseconds]


## Listen to channel
LI (Listen) == 0x94

0x94, channel_no (2 bytes), callsign/id (utf8 string), 0x00, [ pass_token, 0x00 ]

## Carrier key return
Send from a hub to provide carrier key after listening

0x95, channel_no (2 byte), carrier_key (2 bytes)

## Unlisten
UL (unlisten) = 0x96

0x96, carrier_key (4 bytes)


### Key-value pairs
KV (Key Value) == 0x81:

0x80, channel_no (2 bytes), carrier_key (2 bytes), key(utf-8 string), 0x00, value(utf-8 string), 0x00

Used key-value pairs:
DE callsign of sender
IN general info


### carrier events

BE (Bit Event: timed) == 0x82:

0x8
2, channel_key (2 bytes), carrier_key (2 bytes), 8-byte-timestamp (ns), event_type [, ...]

event_type is:
0x00: bit off
0x01: bit on

flags: bitwise flags
Not currently used

## Stream Semantics
0x94 listen (port) <-- listening at source IP:port
# carrier key is return with 0x95
0x81 sending
0x82 start stream -- sets time zero for (fromip, fromport)
0x83 end stream  -- stop keeping time

## Example converstation 

1. Establish a listener socket and send LI (eg. port 0x4001)

0x94, 0x00, 0x40, G, 0, W, C, Z, 0x00

2. hub returns channel and carrier key: 0x95 0x00, 0x40, 0x00, 0x01

3. hub stats sending carriers for the channel eg. 


3. Send something

This sends a dit at about 25wpm
0x82, 0x00, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x05, 0x0, 0x00
key   flags ^---------timestamp(ns)--------------^   on   ^---------timestamp(ns)--------------^   off

## Questions

How about direct udp-udp connections. Why not?


