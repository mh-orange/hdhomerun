package hdhomerun

import (
	"encoding/hex"
	"testing"
)

func TestPacketDump(t *testing.T) {
	expected := "  Type: Discover Request\n"
	expected += "Length: 12\n"
	expected += "  Tags:\n"
	expected += "             Device Type: tuner\n"
	expected += "               Device ID: *\n"
	received := discoverReq.p.dump()
	if received != expected {
		t.Errorf("Unexpected output:\nExpected:\n%s\n     Got:\n%s", expected, received)
	}
}

func TestTlvDumpDeviceType(t *testing.T) {
	tests := []struct {
		expected string
		tag      tlv
	}{
		{
			expected: "       Device Type: *",
			tag: tlv{
				tag:    TagDeviceType,
				length: 0x04,
				value:  DeviceTypeWildcard,
			},
		},
		{
			expected: "       Device Type: tuner",
			tag: tlv{
				tag:    TagDeviceType,
				length: 0x04,
				value:  DeviceTypeTuner,
			},
		},
		{
			expected: "       Device Type: storage",
			tag: tlv{
				tag:    TagDeviceType,
				length: 0x04,
				value:  DeviceTypeStorage,
			},
		},
	}
	for _, test := range tests {
		received := test.tag.dump()
		if test.expected != received {
			t.Errorf("Unexpected output:\nExpected: \"%s\"\n     Got: \"%s\"", test.expected, received)
		}
	}
}

func TestTlvDumpDeviceId(t *testing.T) {
	expected := "         Device ID: 0x01020304"
	received := discoverRpy.p.tags[1].dump()
	if received != expected {
		t.Errorf("Unexpected output:\nExpected: \"%s\"\n     Got: \"%s\"", expected, received)
	}
}
