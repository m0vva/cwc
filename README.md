# bitoip

## What

bit-over-ip: Low level protocol for sending timed on-off state over TCP/IP or UDP/IP

## Why

I'm looking for a way to transmit hand-keyed and keyer-keyed morse over the
internet.  Ideally, this is via a toggled line on a serial port or software
generated from morse keying or morse detecting software.

This becomes generated *carrier* that is a bunch of on and off events, with strict timing maintained so it can be re-constructed at
the far end.

Carriers are grouped into a *channel* which has a name, and some control interface.

Channels are hosted on a *server*.

##  Entities in detail

Here is some more detail on what the logical parts are:

A **carrier** is a set of bit on/off events for a single bit.  A carrier has:
 - an id (which is used to differentiate different bits in a stream)
 - a series of on and off events and control messages
 - some sort of source id (like a string callsign)

Note that a carrier should always end up off, not stuck on.  That means the way we collect and send bit data has to
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

The protocol is composed of packets containing timed on and off information for a bit, with the
necessary time offsets from the start of the stream.

There is a facility to communicate basic key/value ascii pairs.

Generally, top bit set means that is a control value.

multi-byte values are big-endian

## Listen
LI (Listen) == 0x80

0x80, port_number, 0x00 (off) | 0x01 (on)

### Key-value pairs

KV (Key Value) == 0x81:

0x81, length (bytes), bytes, length, value_ascii...

### Bit/stream events

BE (Bit Event: timed) == 0x82:

0x82, flags(1 byte), 4-byte-timestamp (ms) from send start, event_type

event_type is:
0x00: bit off
0x01: bit on

flags: bitwise flags
Not currently used

timestamp=0 for converstation start start

## Stream Semantics
0x80 listen (port) <-- listening at source IP:port
0x82 start stream -- sets time zero for (fromip, fromport)
0x83 end stream  -- stop keeping time

## Example converstation 

1. Establish a listener socket and send LI (eg. port 0x4001)

0x80, 0x40, 0x01, 0x01

2. Set callsign KV pair (optional)

0x81, 0x02, D, E, 0x06, G, 0, W, C, Z

3. Send something

This sends a dit at about 25wpm
0x82, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01
0x82 0x00, 0x00, 0x00, 0x00, 0x30, 0x00








## Questions

How about direct udp-udp connections. Why not?


