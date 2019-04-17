package cwc

import (
	"context"
	"log"
)
import "../bitoip"

// General client
// Can be in CQ mode, in which case all is local muticast on the local network
// Else the client of a reflector
// CQ mode is really simple. Only really have to tx and rx carrier events

func StationClient(ctx context.Context, cqMode bool, addr string, morseIO IO) {

	// if talking to reflector, do some setup
	if !cqMode {



	}
	toSend := make(chan bitoip.CarrierEventPayload)

	go RunRx(ctx, morseIO, toSend)

	for {
		select {
		case <- ctx.Done():
			return
		case cep := <- toSend:
			log.Printf("carrier event payload to send: %v", cep)
		}
	}



	// Reflector mode
	// opt: time sync with server
	// opt: set callsign
	// list channels
	// suscribe channel(s)
	// save carrier id


}