package main

import (
	"flag"
	"net"
	"log"
	"fmt"
	"context"
)

const maxBufferSize = 508

func main() {
	var address = flag.String("address", "localhost:7388", "-address=host:port")

	flag.Parse()

	ReflectorServer(context.TODO(), *address)
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