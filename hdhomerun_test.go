package hdhomerun

type testPacket struct {
	p *Packet
	b []byte
}

// Discover Request
var discoverReq = testPacket{
	p: NewPacket(TypeDiscoverReq, map[TagType]TagValue{
		TagDeviceType: DeviceTypeTuner,
		TagDeviceId:   DeviceIdWildcard,
	}),
	b: []byte{
		0x00, 0x02, 0x00, 0x0c, 0x01, 0x04,
		0x00, 0x00, 0x00, 0x01, 0x02, 0x04,
		0xff, 0xff, 0xff, 0xff, 0x4e, 0x50,
		0x7f, 0x35,
	},
}

// Discover Reply
var discoverRpy = testPacket{
	p: NewPacket(TypeDiscoverRpy, map[TagType]TagValue{
		TagDeviceType: DeviceTypeTuner,
		TagDeviceId:   TagValue{0x01, 0x02, 0x03, 0x04},
	}),
	b: []byte{
		0x00, 0x03, 0x00, 0x0c, 0x01, 0x04,
		0x00, 0x00, 0x00, 0x01, 0x02, 0x04,
		0x01, 0x02, 0x03, 0x04, 0x94, 0x8f,
		0x47, 0xc5,
	},
}

var getRequest = testPacket{}
var getReply = testPacket{}
var setRequest = testPacket{}
var setReply = testPacket{}

var testPackets = []testPacket{
	discoverReq,
	discoverRpy,
}
