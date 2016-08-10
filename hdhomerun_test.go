package hdhomerun

type testPacket struct {
	p *packet
	b []byte
}

// Discover Request
var discoverReq = testPacket{
	p: &packet{
		pktType: TypeDiscoverReq,
		length:  0x0c,
		tags: []tlv{
			tlv{
				tag:    TagDeviceType,
				length: 0x04,
				value:  DeviceTypeTuner,
			},
			tlv{
				tag:    TagDeviceId,
				length: 0x04,
				value:  DeviceIdWildcard,
			},
		},
	},
	b: []byte{
		0x00, 0x02, 0x00, 0x0c, 0x01, 0x04,
		0x00, 0x00, 0x00, 0x01, 0x02, 0x04,
		0xff, 0xff, 0xff, 0xff, 0x4e, 0x50,
		0x7f, 0x35,
	},
}

// Discover Reply
var discoverRpy = testPacket{
	p: &packet{
		pktType: TypeDiscoverRpy,
		length:  0x0c,
		tags: []tlv{
			tlv{
				tag:    TagDeviceType,
				length: 0x04,
				value:  DeviceTypeTuner,
			},
			tlv{
				tag:    TagDeviceId,
				length: 0x04,
				value:  []byte{0x01, 0x02, 0x03, 0x04},
			},
		},
	},
	b: []byte{
		0x00, 0x03, 0x00, 0x0c, 0x01, 0x04,
		0x00, 0x00, 0x00, 0x01, 0x02, 0x04,
		0x01, 0x02, 0x03, 0x04, 0x94, 0x8f,
		0x47, 0xc5,
	},
}

var testPackets = []testPacket{
	discoverReq,
	discoverRpy,
}
