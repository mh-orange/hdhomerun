package hdhomerun

import (
	"testing"
)

func TestPacketDump(t *testing.T) {
	tests := []struct {
		p        *Packet
		expected string
	}{
		{
			p: discoverReq.p,
			expected: "  Type: Discover Request\n" +
				"  Tags:\n" +
				"             Device Type: tuner\n" +
				"               Device ID: *\n",
		}, {
			p:        nil,
			expected: "<nil>",
		},
	}

	for _, test := range tests {
		received := test.p.Dump()
		if received != test.expected {
			t.Errorf("Unexpected output:\nExpected:\n%s\n     Got:\n%s", test.expected, received)
		}
	}
}

func TestPacketTypeString(t *testing.T) {
	tests := []struct {
		pt  PacketType
		str string
	}{
		{TypeDiscoverReq, "Discover Request"},
		{TypeDiscoverRpy, "Discover Reply"},
		{TypeGetSetReq, "Get/Set Request"},
		{TypeGetSetRpy, "Get/Set Reply"},
		{TypeUpgradeReq, "Upgrade Request"},
		{TypeUpgradeRpy, "Upgrade Reply"},
		{PacketType(0xffff), "Unknown"},
	}

	for _, test := range tests {
		expected := test.str
		received := test.pt.String()
		if received != expected {
			t.Errorf("Expected %v to be %s but got %s", test.pt, expected, received)
		}
	}
}

func TestTagTypeString(t *testing.T) {
	tests := []struct {
		tt  TagType
		str string
	}{
		{TagDeviceType, "Device Type"},
		{TagDeviceId, "Device ID"},
		{TagGetSetName, "Get/Set Name"},
		{TagGetSetValue, "Get/Set Value"},
		{TagGetSetLockKey, "Get/Set Lock Key"},
		{TagErrorMsg, "Error Msg"},
		{TagTunerCount, "Tuner Count"},
		{TagDeviceAuthBin, "Device Auth Bin"},
		{TagBaseUrl, "Base URL"},
		{TagDeviceAuthStr, "Device Auth String"},
		{TagType(0xff), "Unknown"},
	}

	for _, test := range tests {
		expected := test.str
		received := test.tt.String()
		if received != expected {
			t.Errorf("Expected %v to be %s but got %s", test.tt, expected, received)
		}
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
				Type:  TagDeviceType,
				Value: DeviceTypeWildcard,
			},
		},
		{
			expected: "       Device Type: tuner",
			tag: Tag{
				Type:  TagDeviceType,
				Value: DeviceTypeTuner,
			},
		},
		{
			expected: "       Device Type: storage",
			tag: Tag{
				Type:  TagDeviceType,
				Value: DeviceTypeStorage,
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
	received := discoverRpy.p.Tags[TagDeviceId].Dump()
	if received != expected {
		t.Errorf("Unexpected output:\nExpected: \"%s\"\n     Got: \"%s\"", expected, received)
	}
}
