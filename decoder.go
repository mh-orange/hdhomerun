package hdhomerun

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"runtime/debug"
)

var (
	ErrCrc error = errors.New("Invalid CRC")
)

type Decoder struct {
	r   io.Reader
	err error
}

func NewDecoder(reader io.Reader) *Decoder {
	d := &Decoder{
		r: reader,
	}
	return d
}

func (d *Decoder) Read(b []byte) (int, error) {
	n := 0
	if d.err == nil {
		n, d.err = d.r.Read(b)
		if d.err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read: %v", d.err)
			debug.PrintStack()
		}
	}
	return n, d.err
}

func (d *Decoder) binaryRead(byteOrder binary.ByteOrder, data interface{}) {
	if d.err == nil {
		d.err = binary.Read(d, byteOrder, data)
	}
}

func (d *Decoder) Decode() (p *Packet, err error) {
	incomingCrc := uint32(0)
	length := uint16(0)
	p = &Packet{}

	d.binaryRead(binary.BigEndian, &p.Type)
	d.binaryRead(binary.BigEndian, &length)

	data := make([]byte, length)
	d.Read(data)
	computedCrc := crc32.Update(0, crc32.IEEETable, []byte{byte(p.Type >> 8), byte(p.Type)})
	computedCrc = crc32.Update(computedCrc, crc32.IEEETable, []byte{byte(length >> 8), byte(length)})
	computedCrc = crc32.Update(computedCrc, crc32.IEEETable, data)

	d.binaryRead(binary.LittleEndian, &incomingCrc)
	if d.err == nil && incomingCrc != computedCrc {
		d.err = ErrCrc
	}

	var tagLength uint16
	buffer := bytes.NewReader(data)
	d.r = buffer
	for d.err == nil && buffer.Len() > 0 {
		t := Tag{}
		var lsb, msb uint8

		d.binaryRead(binary.BigEndian, &t.Tag)
		d.binaryRead(binary.BigEndian, &lsb)

		// two byte length
		if lsb&0x80 == 0x80 {
			d.binaryRead(binary.BigEndian, &msb)
			tagLength = uint16(lsb&0x7f) | uint16(msb)<<7
		} else {
			tagLength = uint16(lsb)
		}

		t.Value = make([]byte, tagLength)
		d.Read(t.Value)
		if d.err == nil {
			p.Tags = append(p.Tags, t)
		}
	}
	return p, d.err
}
