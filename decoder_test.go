package hdhomerun

import (
	"bytes"
	"reflect"
	"testing"
)

func decode(b []byte) (*Packet, error) {
	buffer := bytes.NewBuffer(b)
	d := NewDecoder(buffer)
	return d.Next()
}

func TestDecode(t *testing.T) {
	for _, testPacket := range testPackets {
		p, err := decode(testPacket.b)
		if err != nil {
			t.Errorf("Packet decoding failed with error: %v", err)
		}

		if !reflect.DeepEqual(testPacket.p, p) {
			t.Errorf("Packet decoding failed\nExpected:\n%s\nReceived:\n%s\n", testPacket.p.Dump(), p.Dump())
		}

		// fiddle with the crc to generate an error
		b := make([]byte, len(testPacket.b))
		copy(b, testPacket.b)
		b[len(b)-1] = ^b[len(b)-1]
		p, err = decode(b)
		if err != ErrCrc {
			t.Errorf("Expected error %v but got %v", ErrCrc, err)
		}
	}
}
