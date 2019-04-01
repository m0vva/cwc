package main

import (
	"flag"
	"net"
			"context"
	"../bitoip"
)

func main() {
	var address = flag.String("address", "localhost:7388", "-address=host:port")

	flag.Parse()

	ReflectorServer(context.TODO(), *address)
}

func ReflectorServer(ctx context.Context, address string) {
	bitoip.UDPRx(ctx, address, Handler)
}

/**
	Handle an incoming message to the reflector
 */
func Handler(verb bitoip.MessageVerb, payload bitoip.Payload, address net.Addr) {
	switch verb {
	case bitoip.EnumerateChannels:

	}
}