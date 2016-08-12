package hdhomerun

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

var (
	// Not sure if these should be null terminated.  Documentation in
	// libhdhomerun does not specify null terminated strings for
	// discovery packets where the following values are used
	DeviceTypeWildcard []byte = []byte{0xFF, 0xFF, 0xFF, 0xFF}
	DeviceTypeTuner    []byte = []byte{0x00, 0x00, 0x00, 0x01}
	DeviceTypeStorage  []byte = []byte{0x00, 0x00, 0x00, 0x05}
	DeviceIdWildcard   []byte = []byte{0xFF, 0xFF, 0xFF, 0xFF}
)

type PacketType uint16

const (
	TypeDiscoverReq PacketType = 0x0002
	TypeDiscoverRpy PacketType = 0x0003
	TypeGetSetReq   PacketType = 0x0004
	TypeGetSetRpy   PacketType = 0x0005
	TypeUpgradeReq  PacketType = 0x0006
	TypeUpgradeRpy  PacketType = 0x0007
)

func (pt PacketType) String() string {
	switch pt {
	case TypeDiscoverReq:
		return "Discover Request"
	case TypeDiscoverRpy:
		return "Discover Reply"
	case TypeGetSetReq:
		return "Get/Set Request"
	case TypeGetSetRpy:
		return "Get/Set Reply"
	case TypeUpgradeReq:
		return "Upgrade Request"
	case TypeUpgradeRpy:
		return "Upgrade Reply"
	}

	return "Unknown"
}

type Packet struct {
	Type PacketType
	Tags []Tag
}

//func NewPacket(t PacketType, tags map[TagType]TagValue) *Packet {
//}

func (p *Packet) Dump() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("  Type: %s\n", p.Type))
	buffer.WriteString("  Tags:\n")
	for _, t := range p.Tags {
		buffer.WriteString(fmt.Sprintf("      %s\n", t.Dump()))
	}
	return buffer.String()
}

type TagType uint8

const (
	TagDeviceType    TagType = 0x01
	TagDeviceId      TagType = 0x02
	TagGetSetName    TagType = 0x03
	TagGetSetValue   TagType = 0x04
	TagGetSetLockKey TagType = 0x15
	TagErrorMsg      TagType = 0x05
	TagTunerCount    TagType = 0x10
	TagDeviceAuthBin TagType = 0x29
	TagBaseUrl       TagType = 0x2A
	TagDeviceAuthStr TagType = 0x2B
)

func (tt TagType) String() string {
	switch tt {
	case TagDeviceType:
		return "Device Type"
	case TagDeviceId:
		return "Device ID"
	case TagGetSetName:
		return "Get/Set Name"
	case TagGetSetValue:
		return "Get/Set Value"
	case TagGetSetLockKey:
		return "Get/Set Lock Key"
	case TagErrorMsg:
		return "Error Msg"
	case TagTunerCount:
		return "Tuner Count"
	case TagDeviceAuthBin:
		return "Device Auth Bin"
	case TagBaseUrl:
		return "Base URL"
	case TagDeviceAuthStr:
		return "Device Auth String"
	}

	return "Unknown"
}

type TagValue []byte

func (tv TagValue) String() string {
	return string(tv)
}

type Tag struct {
	Tag   TagType
	Value TagValue
}

func (t Tag) String() string {
	return string(t.Value)
}

func (t *Tag) Dump() string {
	value := t.String()
	if t.Tag == TagDeviceType {
		if bytes.Equal(t.Value, DeviceTypeWildcard) {
			value = "*"
		} else if bytes.Equal(t.Value, DeviceTypeTuner) {
			value = "tuner"
		} else if bytes.Equal(t.Value, DeviceTypeStorage) {
			value = "storage"
		}
	} else if t.Tag == TagDeviceId {
		if bytes.Equal(t.Value, DeviceIdWildcard) {
			value = "*"
		} else {
			value = fmt.Sprintf("0x%s", hex.EncodeToString(t.Value))
		}
	}

	return fmt.Sprintf("%18s: %s", t.Tag, value)
}
