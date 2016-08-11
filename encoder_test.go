package hdhomerun

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestEncode(t *testing.T) {
	for _, testPacket := range testPackets {
		buffer := bytes.NewBuffer([]byte{})
		e := NewEncoder(buffer)
		e.Encode(testPacket.p)
		if !bytes.Equal(testPacket.b, buffer.Bytes()) {
			t.Errorf("Packet encoding failed\nExpected:\n%s\nReceived:\n%s\n", hex.Dump(testPacket.b), hex.Dump(buffer.Bytes()))
		}
	}
}
