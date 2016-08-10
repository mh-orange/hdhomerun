package hdhomerun

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	for _, testPacket := range testPackets {
		buffer := bytes.NewBuffer([]byte{})
		buffer.Write(testPacket.b)
		d := newDecoder(buffer)
		p, err := d.decode()
		if err != nil {
			t.Errorf("Packet decoding failed with error: %v", err)
		}

		if !reflect.DeepEqual(testPacket.p, p) {
			t.Errorf("Packet decoding failed\nExpected:\n%s\nReceived:\n%s\n", testPacket.p.Dump(), p.Dump())
		}
	}
}
