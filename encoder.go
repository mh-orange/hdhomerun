package hdhomerun

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
)

type Encoder struct {
	w   io.Writer
	err error
}

func NewEncoder(writer io.Writer) *Encoder {
	e := &Encoder{
		w: writer,
	}
	return e
}

func (e *Encoder) Write(p []byte) (int, error) {
	n := 0
	if e.err == nil {
		n, e.err = e.w.Write(p)
	}
	return n, e.err
}

func (e *Encoder) Encode(p *Packet) error {
	buffer := bytes.NewBuffer(make([]byte, 0))
	length := uint16(0)
	for _, tag := range p.Tags {
		length += 2 + uint16(len(tag.Value))
		if len(tag.Value) > 127 {
			length += 1
		}
	}

	binary.Write(buffer, binary.BigEndian, p.Type)
	binary.Write(buffer, binary.BigEndian, length)
	for _, t := range p.Tags {
		buffer.Write([]byte{byte(t.Type)})
		length := uint8(len(t.Value))
		if length > 127 {
			lsb := 0x80 | (length & 0x00ff)
			msb := length >> 7
			buffer.Write([]byte{byte(lsb), byte(msb)})
		} else {
			buffer.Write([]byte{byte(length)})
		}

		buffer.Write(t.Value)
	}

	crc := crc32.ChecksumIEEE(buffer.Bytes())
	binary.Write(buffer, binary.LittleEndian, crc)
	buffer.WriteTo(e)
	return e.err
}
