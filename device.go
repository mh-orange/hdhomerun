package hdhomerun

import (
	"bytes"
	"encoding/hex"
	"io"
	"net"
	"time"
)

type Device struct {
	Connection
	id   []byte
	e    *Encoder
	d    *Decoder
	Addr string
}

type discoverConn interface {
	Close() error
	ReadFrom([]byte) (int, net.Addr, error)
	SetReadDeadline(time.Time) error
	WriteTo([]byte, net.Addr) (int, error)
}

var listenUDP = func() (discoverConn, error) {
	return net.ListenUDP("udp", nil)
}

func Discover(ip net.IP, timeout time.Duration) (chan *Device, error) {
	if ip == nil {
		ip = net.IPv4bcast
	}

	conn, err := listenUDP()
	if err != nil {
		return nil, err
	}

	devices := make(chan *Device, 10)
	go func() {
		buffer := make([]byte, 1500)
		readBuffer := bytes.NewBuffer(make([]byte, 1500))
		readBuffer.Reset()
		decoder := NewDecoder(readBuffer)
		discovered := make(map[string]bool)
		defer close(devices)

		end := time.Now().Add(timeout)
		for t := time.Now(); t.Before(end); t = time.Now() {
			conn.SetReadDeadline(time.Now().Add(timeout))

			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					continue
				} else if err != io.EOF {
					Logger.Printf("1) %v", err)
				}
				break
			}

			readBuffer.Write(buffer[0:n])
			p, err := decoder.Next()

			if err != nil {
				Logger.Printf("%v", err)
				break
			}

			if p.Type != TypeDiscoverRpy {
				continue
			}

			if _, found := discovered[addr.String()]; !found {
				devices <- NewDevice(NewTCPConnection(addr.String()), p.Tags[TagDeviceId].Value, addr.String())
			}
			discovered[addr.String()] = true
		}
		conn.Close()
	}()

	writeBuffer := bytes.NewBuffer([]byte{})
	encoder := NewEncoder(writeBuffer)
	go func() {
		for i := 0; i < 2; i++ {
			err := encoder.Encode(NewPacket(TypeDiscoverReq, map[TagType]TagValue{
				TagDeviceType: DeviceTypeTuner,
				TagDeviceId:   DeviceIdWildcard,
			}))
			if err == nil {
				_, err = conn.WriteTo(writeBuffer.Bytes(), &net.UDPAddr{IP: ip, Port: 65001})
			}
			if err != nil {
				Logger.Printf("Failed to send discovery packet: %v", err)
			}
		}
	}()

	return devices, nil
}

func NewDevice(conn Connection, id []byte, addr string) *Device {
	return &Device{
		Connection: conn,
		id:         id,
		Addr:       addr,
	}
}

func (d *Device) ID() string {
	return hex.EncodeToString(d.id)
}

func (d *Device) getset(name string, value *string) (resp string, err error) {
	tags := make(map[TagType]TagValue)
	tags[TagGetSetName] = TagValue(name)

	if value != nil {
		tags[TagGetSetValue] = TagValue(*value)
	}

	err = d.Send(NewPacket(TypeGetSetReq, tags))
	if err == nil {
		var p *Packet
		p, err = d.Recv()
		if p.Type != TypeGetSetRpy {
			err = wrongPacketType(TypeGetSetRpy, p.Type)
		} else {
			resp = p.Tags[TagGetSetValue].String()
			if tag, found := p.Tags[TagErrorMsg]; found {
				err = ErrRemoteError(tag.String())
			}
		}
	}
	return
}

func (d *Device) Get(name string) (string, error) {
	return d.getset(name, nil)
}

func (d *Device) Set(name, value string) (string, error) {
	return d.getset(name, &value)
}

func (d *Device) Tuner(n int) *Tuner {
	return newTuner(d, n)
}
