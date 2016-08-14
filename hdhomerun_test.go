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

var getReq = testPacket{
	p: NewPacket(TypeGetSetReq, map[TagType]TagValue{
		TagGetSetName: TagValue("/card/status"),
	}),
	b: []byte{
		0x00, 0x04, 0x00, 0x0f, 0x03, 0x0d,
		0x2f, 0x63, 0x61, 0x72, 0x64, 0x2f,
		0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
		0x00, 0x9f, 0x93, 0x24, 0x2b,
	},
}

var getRpy = testPacket{
	p: NewPacket(TypeGetSetRpy, map[TagType]TagValue{
		TagGetSetName:  TagValue("help"),
		TagGetSetValue: TagValue("Supported configuration options:\n/card/status\n/ir/target <protocol>://<ip>:<port>\n/lineup/scan\n/oob/channel <modulation>:<freq>\n/oob/debug\n/oob/status\n/sys/copyright\n/sys/debug\n/sys/features\n/sys/hwmodel\n/sys/model\n/sys/restart <resource>\n/sys/version\n/tuner<n>/channel <modulation>:<freq|ch>\n/tuner<n>/channelmap <channelmap>\n/tuner<n>/debug\n/tuner<n>/filter \"0x<nnnn>-0x<nnnn> [...]\"\n/tuner<n>/lockkey\n/tuner<n>/program <program number>\n/tuner<n>/streaminfo\n/tuner<n>/status\n/tuner<n>/target <ip>:<port>\n/tuner<n>/vchannel <vchannel>\n/tuner<n>/vstatus\n"),
	}),
	b: []byte{
		0x00, 0x05, 0x02, 0x35, 0x03, 0x05, 0x68, 0x65, 0x6c, 0x70, 0x00, 0x04,
		0xab, 0x04, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x20,
		0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
		0x6e, 0x20, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x3a, 0x0a, 0x2f,
		0x63, 0x61, 0x72, 0x64, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x0a,
		0x2f, 0x69, 0x72, 0x2f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x20, 0x3c,
		0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x3e, 0x3a, 0x2f, 0x2f,
		0x3c, 0x69, 0x70, 0x3e, 0x3a, 0x3c, 0x70, 0x6f, 0x72, 0x74, 0x3e, 0x0a,
		0x2f, 0x6c, 0x69, 0x6e, 0x65, 0x75, 0x70, 0x2f, 0x73, 0x63, 0x61, 0x6e,
		0x0a, 0x2f, 0x6f, 0x6f, 0x62, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65,
		0x6c, 0x20, 0x3c, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f,
		0x6e, 0x3e, 0x3a, 0x3c, 0x66, 0x72, 0x65, 0x71, 0x3e, 0x0a, 0x2f, 0x6f,
		0x6f, 0x62, 0x2f, 0x64, 0x65, 0x62, 0x75, 0x67, 0x0a, 0x2f, 0x6f, 0x6f,
		0x62, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x0a, 0x2f, 0x73, 0x79,
		0x73, 0x2f, 0x63, 0x6f, 0x70, 0x79, 0x72, 0x69, 0x67, 0x68, 0x74, 0x0a,
		0x2f, 0x73, 0x79, 0x73, 0x2f, 0x64, 0x65, 0x62, 0x75, 0x67, 0x0a, 0x2f,
		0x73, 0x79, 0x73, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73,
		0x0a, 0x2f, 0x73, 0x79, 0x73, 0x2f, 0x68, 0x77, 0x6d, 0x6f, 0x64, 0x65,
		0x6c, 0x0a, 0x2f, 0x73, 0x79, 0x73, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
		0x0a, 0x2f, 0x73, 0x79, 0x73, 0x2f, 0x72, 0x65, 0x73, 0x74, 0x61, 0x72,
		0x74, 0x20, 0x3c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x3e,
		0x0a, 0x2f, 0x73, 0x79, 0x73, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
		0x6e, 0x0a, 0x2f, 0x74, 0x75, 0x6e, 0x65, 0x72, 0x3c, 0x6e, 0x3e, 0x2f,
		0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x20, 0x3c, 0x6d, 0x6f, 0x64,
		0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3e, 0x3a, 0x3c, 0x66, 0x72,
		0x65, 0x71, 0x7c, 0x63, 0x68, 0x3e, 0x0a, 0x2f, 0x74, 0x75, 0x6e, 0x65,
		0x72, 0x3c, 0x6e, 0x3e, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
		0x6d, 0x61, 0x70, 0x20, 0x3c, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
		0x6d, 0x61, 0x70, 0x3e, 0x0a, 0x2f, 0x74, 0x75, 0x6e, 0x65, 0x72, 0x3c,
		0x6e, 0x3e, 0x2f, 0x64, 0x65, 0x62, 0x75, 0x67, 0x0a, 0x2f, 0x74, 0x75,
		0x6e, 0x65, 0x72, 0x3c, 0x6e, 0x3e, 0x2f, 0x66, 0x69, 0x6c, 0x74, 0x65,
		0x72, 0x20, 0x22, 0x30, 0x78, 0x3c, 0x6e, 0x6e, 0x6e, 0x6e, 0x3e, 0x2d,
		0x30, 0x78, 0x3c, 0x6e, 0x6e, 0x6e, 0x6e, 0x3e, 0x20, 0x5b, 0x2e, 0x2e,
		0x2e, 0x5d, 0x22, 0x0a, 0x2f, 0x74, 0x75, 0x6e, 0x65, 0x72, 0x3c, 0x6e,
		0x3e, 0x2f, 0x6c, 0x6f, 0x63, 0x6b, 0x6b, 0x65, 0x79, 0x0a, 0x2f, 0x74,
		0x75, 0x6e, 0x65, 0x72, 0x3c, 0x6e, 0x3e, 0x2f, 0x70, 0x72, 0x6f, 0x67,
		0x72, 0x61, 0x6d, 0x20, 0x3c, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d,
		0x20, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x3e, 0x0a, 0x2f, 0x74, 0x75,
		0x6e, 0x65, 0x72, 0x3c, 0x6e, 0x3e, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x61,
		0x6d, 0x69, 0x6e, 0x66, 0x6f, 0x0a, 0x2f, 0x74, 0x75, 0x6e, 0x65, 0x72,
		0x3c, 0x6e, 0x3e, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x0a, 0x2f,
		0x74, 0x75, 0x6e, 0x65, 0x72, 0x3c, 0x6e, 0x3e, 0x2f, 0x74, 0x61, 0x72,
		0x67, 0x65, 0x74, 0x20, 0x3c, 0x69, 0x70, 0x3e, 0x3a, 0x3c, 0x70, 0x6f,
		0x72, 0x74, 0x3e, 0x0a, 0x2f, 0x74, 0x75, 0x6e, 0x65, 0x72, 0x3c, 0x6e,
		0x3e, 0x2f, 0x76, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x20, 0x3c,
		0x76, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x3e, 0x0a, 0x2f, 0x74,
		0x75, 0x6e, 0x65, 0x72, 0x3c, 0x6e, 0x3e, 0x2f, 0x76, 0x73, 0x74, 0x61,
		0x74, 0x75, 0x73, 0x0a, 0x00, 0x25, 0x81, 0x34, 0x87,
	},
}

var setRequest = testPacket{}
var setReply = testPacket{}

var testPackets = []testPacket{
	discoverReq,
	discoverRpy,
	getReq,
	getRpy,
}
