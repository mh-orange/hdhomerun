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

const (
	TypeDiscoverReq uint16 = 0x0002
	TypeDiscoverRpy uint16 = 0x0003
	TypeGetSetReq   uint16 = 0x0004
	TypeGetSetRpy   uint16 = 0x0005
	TypeUpgradeReq  uint16 = 0x0006
	TypeUpgradeRpy  uint16 = 0x0007

	TagDeviceType    uint8 = 0x01
	TagDeviceId      uint8 = 0x02
	TagGetSetName    uint8 = 0x03
	TagGetSetValue   uint8 = 0x04
	TagGetSetLockKey uint8 = 0x15
	TagErrorMsg      uint8 = 0x05
	TagTunerCount    uint8 = 0x10
	TagDeviceAuthBin uint8 = 0x29
	TagBaseUrl       uint8 = 0x2A
	TagDeviceAuthStr uint8 = 0x2B
)

var (
	typeNames = map[uint16]string{
		TypeDiscoverReq: "Discover Request",
		TypeDiscoverRpy: "Discover Reply",
		TypeGetSetReq:   "Get/Set Request",
		TypeGetSetRpy:   "Get/Set Reply",
		TypeUpgradeReq:  "Upgrade Request",
		TypeUpgradeRpy:  "Upgrade Reply",
	}

	tagNames = map[uint8]string{
		TagDeviceType:    "Device Type",
		TagDeviceId:      "Device ID",
		TagGetSetName:    "Get/Set Name",
		TagGetSetValue:   "Get/Set Value",
		TagGetSetLockKey: "Get/Set Lock Key",
		TagErrorMsg:      "Error Msg",
		TagTunerCount:    "Tuner Count",
		TagDeviceAuthBin: "Device Auth Bin",
		TagBaseUrl:       "Base URL",
		TagDeviceAuthStr: "Device Auth String",
	}
)

type Packet struct {
	Type uint16
	Tags []Tag
}

func (p *Packet) Dump() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("  Type: %s\n", typeNames[p.Type]))
	buffer.WriteString("  Tags:\n")
	for _, t := range p.Tags {
		buffer.WriteString(fmt.Sprintf("      %s\n", t.Dump()))
	}
	return buffer.String()
}

type Tag struct {
	Tag   uint8
	Value []byte
}

func (t *Tag) String() string {
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

	return fmt.Sprintf("%18s: %s", tagNames[t.Tag], value)
}
