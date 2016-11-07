# bitoip

## What

bit-over-ip: Low level protocol for sending timed on-off state over TCP/IP or UDP/IP

## Why

I'm looking for a way to transmit hand-keyed and keyer-keyed morse over the
internet.  Ideally, this is via one line on a serial port or even software generated.

There is one bit that is toggled on and off, and the timing is important, so something that 
is stream-oriented in that it can reconstruct time is important. The bandwidth needed
is some small number of bits per second.

## Protocol

The protocol is composed of packets containing on and off information 
and time offsets from the start of connection or absolute timestamps.  

There is a facility to communicate basic key/value ascii pairs.

Generally, top bit set means that is a control value.

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


