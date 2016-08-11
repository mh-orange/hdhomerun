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

type decoder struct {
	r   io.Reader
	err error
}

func newDecoder(reader io.Reader) *decoder {
	d := &decoder{
		r: reader,
	}
	return d
}

func (d *decoder) Read(b []byte) (int, error) {
	n := 0
	if d.err == nil {
		n, d.err = d.r.Read(b)
	}
	return n, d.err
}

func (d *decoder) binaryRead(byteOrder binary.ByteOrder, data interface{}) {
	if d.err == nil {
		d.err = binary.Read(d, byteOrder, data)
	}
}

func (d *decoder) decode() (p *packet, err error) {
	incomingCrc := uint32(0)
	p = &packet{}

	d.binaryRead(binary.BigEndian, &p.pktType)
	d.binaryRead(binary.BigEndian, &p.length)

	data := make([]byte, p.length)
	d.Read(data)
	computedCrc := crc32.Update(0, crc32.IEEETable, []byte{byte(p.pktType >> 8), byte(p.pktType)})
	computedCrc = crc32.Update(computedCrc, crc32.IEEETable, []byte{byte(p.length >> 8), byte(p.length)})
	computedCrc = crc32.Update(computedCrc, crc32.IEEETable, data)

	d.binaryRead(binary.LittleEndian, &incomingCrc)
	if d.err == nil && incomingCrc != computedCrc {
		d.err = ErrCrc
	}

	buffer := bytes.NewReader(data)
	d.r = buffer
	for d.err == nil && buffer.Len() > 0 {
		t := tlv{}
		var lsb, msb uint8

		d.binaryRead(binary.BigEndian, &t.tag)
		d.binaryRead(binary.BigEndian, &lsb)

		// two byte length
		if lsb&0x80 == 0x80 {
			d.binaryRead(binary.BigEndian, &msb)
			t.length = uint16(lsb&0x7f) | uint16(msb)<<7
		} else {
			t.length = uint16(lsb)
		}

		t.value = make([]byte, t.length)
		d.Read(t.value)
		if d.err == nil {
			p.tags = append(p.tags, t)
		}
	}
	return p, d.err
}
