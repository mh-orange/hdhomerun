package hdhomerun

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
)

type encoder struct {
	w   io.Writer
	err error
}

func newEncoder(writer io.Writer) *encoder {
	e := &encoder{
		w: writer,
	}
	return e
}

func (e *encoder) Write(p []byte) (int, error) {
	if e.err == nil {
		return e.w.Write(p)
	}
	return 0, e.err
}

func (e *encoder) encode(p *packet) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	for _, t := range p.tags {
		buffer.Write([]byte{byte(t.tag)})
		if t.length > 127 {
			lsb := 0x80 | (t.length & 0x00ff)
			msb := t.length >> 7
			buffer.Write([]byte{byte(lsb), byte(msb)})
		} else {
			buffer.Write([]byte{byte(t.length)})
		}

		buffer.Write(t.value)
	}

	crc := crc32.ChecksumIEEE(buffer.Bytes())
	binary.Write(e, binary.BigEndian, p.pktType)
	binary.Write(e, binary.BigEndian, p.length)
	buffer.WriteTo(e)
	binary.Write(e, binary.LittleEndian, crc)
}
