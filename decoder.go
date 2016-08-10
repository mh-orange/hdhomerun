package hdhomerun

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
)

var (
	ErrCrc error = errors.New("Invalid CRC")
)

type decoder struct {
	r    io.Reader
	data []byte
	crc  uint32
	err  error
}

type decoderState func(*packet) decoderState

func newDecoder(reader io.Reader) *decoder {
	d := &decoder{
		r: reader,
	}
	return d
}

func (d *decoder) decode() (p *packet, err error) {
	p = &packet{}
	for state := d.readPktType; state != nil && d.err == nil; {
		state = state(p)
	}
	return
}

func (d *decoder) readPktType(p *packet) decoderState {
	d.err = binary.Read(d.r, binary.BigEndian, &p.pktType)
	return d.readLength
}

func (d *decoder) readLength(p *packet) decoderState {
	d.err = binary.Read(d.r, binary.BigEndian, &p.length)
	return d.readData
}

func (d *decoder) readData(p *packet) decoderState {
	d.data = make([]byte, p.length)
	_, d.err = d.r.Read(d.data)
	if d.err == nil {
		d.crc = crc32.ChecksumIEEE(d.data)
	}
	return d.readCrc
}

func (d *decoder) readCrc(p *packet) decoderState {
	crc := uint32(0)
	d.err = binary.Read(d.r, binary.LittleEndian, &crc)
	if d.err == nil {
		if d.crc != crc {
			d.err = ErrCrc
		}
	}
	return d.parseTags
}

func (d *decoder) parseTags(p *packet) decoderState {
	for len(d.data) > 0 {
		t := tlv{}
		if len(d.data) < 2 {
			d.err = fmt.Errorf("Failed to parse tag: missing tag or length")
			return nil
		}

		t.tag = uint8(d.data[0])
		d.data = d.data[1:]
		// two byte length
		if d.data[0]&0x80 == 0x80 {
			if len(d.data) < 1 {
				d.err = fmt.Errorf("Failed to parse tag: varlen length but no high order byte")
				return nil
			}

			t.length = uint16(d.data[0] & 0x7f)
			d.data = d.data[1:]
			t.length |= uint16(d.data[0]) << 7
		} else {
			t.length = uint16(d.data[0])
		}
		d.data = d.data[1:]
		if len(d.data) < int(t.length) {
			d.err = fmt.Errorf("Failed to parse tag: remaining data less than tag length")
			return nil
		}
		t.value = d.data[0:t.length]
		d.data = d.data[t.length:]
		p.tags = append(p.tags, t)
	}
	return nil
}
