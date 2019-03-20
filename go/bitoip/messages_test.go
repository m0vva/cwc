package bitoip

import (
	"testing"
	"fmt"
)

func TestEncodePayload(t *testing.T) {
	verb := TimeSync
	payload := TimeSyncPayload{1}

	myBytes := EncodePayload(verb, payload)
	fmt.Println( myBytes);
}
