package hdhomerun

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

var (
	ErrCrc error = errors.New("Invalid CRC")
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{
		r: reader,
	}
}

func (d *Decoder) Next() (*Packet, error) {
	pd := &packetDecoder{
		r: d.r,
	}
	return pd.decode()
}

type packetDecoder struct {
	r   io.Reader
	err error
}

func (pd *packetDecoder) Read(b []byte) (int, error) {
	n := 0
	if pd.err == nil {
		n, pd.err = pd.r.Read(b)
	}
	return n, pd.err
}

func (pd *packetDecoder) binaryRead(byteOrder binary.ByteOrder, data interface{}) {
	if pd.err == nil {
		pd.err = binary.Read(pd, byteOrder, data)
	}
}

func (pd *packetDecoder) decode() (p *Packet, err error) {
	incomingCrc := uint32(0)
	length := uint16(0)
	p = &Packet{
		Tags: make(map[TagType]Tag),
	}

	pd.binaryRead(binary.BigEndian, &p.Type)
	pd.binaryRead(binary.BigEndian, &length)

	data := make([]byte, length)
	pd.Read(data)
	computedCrc := crc32.Update(0, crc32.IEEETable, []byte{byte(p.Type >> 8), byte(p.Type)})
	computedCrc = crc32.Update(computedCrc, crc32.IEEETable, []byte{byte(length >> 8), byte(length)})
	computedCrc = crc32.Update(computedCrc, crc32.IEEETable, data)

	pd.binaryRead(binary.LittleEndian, &incomingCrc)
	if pd.err == nil && incomingCrc != computedCrc {
		pd.err = ErrCrc
	}

	var tagLength uint16
	buffer := bytes.NewReader(data)
	pd.r = buffer
	for pd.err == nil && buffer.Len() > 0 {
		t := Tag{}
		var lsb, msb uint8

		pd.binaryRead(binary.BigEndian, &t.Type)
		pd.binaryRead(binary.BigEndian, &lsb)

		// two byte length
		if lsb&0x80 == 0x80 {
			pd.binaryRead(binary.BigEndian, &msb)
			tagLength = uint16(lsb&0x7f) | uint16(msb)<<7
		} else {
			tagLength = uint16(lsb)
		}

		t.Value = make([]byte, tagLength)
		pd.Read(t.Value)
		if pd.err == nil {
			p.Tags[t.Type] = t
		}
	}
	return p, pd.err
}
