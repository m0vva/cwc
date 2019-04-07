package main

import (
	"flag"
	"net"
	"log"
	"fmt"
	"context"
	"../cwc"
)

const maxBufferSize = 508
const CqPort = 5990
const LocalMulticast = "224.0.0.73:%d"

func main() {
	var refAddPtr = flag.String("ref", "cwc,nodestone.io:7388", "--ref=host:port")
	var cqPtr = flag.Bool("cq", false, "--cq is CQ mode, no server, local broadcast")
	var localPort = flag.Int("port", CqPort, "--port=<local-udp-port>")

	flag.Parse()

	cqMode := *cqPtr
	refAddress := *refAddPtr


	if (cqMode) {
		mcAddress := fmt.Sprintf(LocalMulticast, *localPort)
		log.Printf("Starting in CQ mode with local multicast address %s", mcAddress)

		StationClient(true, mcAddress)
	} else {
		log.Printf("Connecting to reflector %s", refAddress)

		StationClient(false, refAddress)
	}
}

// General client
// Can be in CQ mode, in which case all is local muticast on the local network
// Else the client of a reflector
// CQ mode is really simple. Only really have to tx and rx carrier events
func StationClient(cqMode bool, addr string) {
	// CQ mode
	// listen on mc address
	// look for bit events and send them
	// send using any channel
	// rx all channels
	// that's it
	cwc.RunRx()

	// Reflector mode
	// opt: time sync with server
	// opt: set callsign
	// list channels
	// suscribe channel(s)
	// save carrier id


}


func ReflectorServer(ctx context.Context, address string) {
	pc, err := net.ListenPacket("udp", address)

	if err != nil {
		return
	}

	defer pc.Close()

	buffer := make([]byte, maxBufferSize)
	doneChan := make(chan error, 1)

	go func() {
		for {
			_, _, err := pc.ReadFrom(buffer)

			if err != nil {
				doneChan <- err
				return
			}
			log.Printf("packet rx: %#v", buffer)
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

}