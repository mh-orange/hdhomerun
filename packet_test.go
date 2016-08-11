package hdhomerun

import (
	"testing"
)

func TestPacketDump(t *testing.T) {
	expected := "  Type: Discover Request\n"
	expected += "  Tags:\n"
	expected += "             Device Type: tuner\n"
	expected += "               Device ID: *\n"
	received := discoverReq.p.Dump()
	if received != expected {
		t.Errorf("Unexpected output:\nExpected:\n%s\n     Got:\n%s", expected, received)
	}
}

func TestTlvDumpDeviceType(t *testing.T) {
	tests := []struct {
		expected string
		tag      Tag
	}{
		{
			expected: "       Device Type: *",
			tag: Tag{
				tag:   TagDeviceType,
				value: DeviceTypeWildcard,
			},
		},
		{
			expected: "       Device Type: tuner",
			tag: Tag{
				tag:   TagDeviceType,
				value: DeviceTypeTuner,
			},
		},
		{
			expected: "       Device Type: storage",
			tag: Tag{
				tag:   TagDeviceType,
				value: DeviceTypeStorage,
			},
		},
	}
	for _, test := range tests {
		received := test.tag.Dump()
		if test.expected != received {
			t.Errorf("Unexpected output:\nExpected: \"%s\"\n     Got: \"%s\"", test.expected, received)
		}
	}
}

func TestTlvDumpDeviceId(t *testing.T) {
	expected := "         Device ID: 0x01020304"
	received := discoverRpy.p.Tags[1].Dump()
	if received != expected {
		t.Errorf("Unexpected output:\nExpected: \"%s\"\n     Got: \"%s\"", expected, received)
	}
}
