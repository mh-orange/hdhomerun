package hdhomerun

type testPacket struct {
	p *packet
	b []byte
}

var testPackets = []testPacket{
	// Discover Request
	testPacket{
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
			0xff, 0xff, 0xff, 0xff, 0x1e, 0xd5,
			0x2d, 0x90,
		},
	},
	// Discover Reply
	testPacket{
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
			0x01, 0x02, 0x03, 0x04, 0x2c, 0xd1,
			0xee, 0xd9,
		},
	},
}
