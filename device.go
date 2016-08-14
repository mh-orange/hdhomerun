package hdhomerun

import (
	"bytes"
	"encoding/hex"
	"net"
	"time"
)

type Device struct {
	Connection
	id []byte
	e  *Encoder
	d  *Decoder
}

func Discover(ip net.IP, timeout time.Duration) (map[string]*Device, error) {
	if ip == nil {
		ip = net.IPv4bcast
	}
	devices := make(chan *Device, 10)

	laddr, _ := net.ResolveUDPAddr("udp", ":0")
	conn, err := net.ListenUDP("udp", laddr)
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(devices)
		var raddr *net.UDPAddr
		var n int
		var p *Packet

		buffer := make([]byte, 1024)
		end := time.Now().Add(timeout)
		for t := time.Now(); t.Before(end); t = time.Now() {
			conn.SetReadDeadline(time.Now().Add((10 * time.Millisecond)))
			n, raddr, err = conn.ReadFromUDP(buffer)
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					err = nil
					continue
				}
				break
			}
			decoder := NewDecoder(bytes.NewBuffer(buffer[0:n]))
			p, err = decoder.Next()
			if err != nil {
				break
			}
			if p.Type != TypeDiscoverRpy {
				continue
			}
			device := NewDevice(NewTcpConnection(raddr.IP, raddr.Port), p.Tags[TagDeviceId].Value)
			devices <- device
		}
	}()

	discoveredDevices := make(map[string]*Device)
	for i := 0; i < 2; i++ {
		writeBuffer := bytes.NewBuffer([]byte{})
		encoder := NewEncoder(writeBuffer)
		encoder.Encode(NewPacket(TypeDiscoverReq, map[TagType]TagValue{
			TagDeviceType: DeviceTypeTuner,
			TagDeviceId:   DeviceIdWildcard,
		}))
		_, err = conn.WriteTo(writeBuffer.Bytes(), &net.UDPAddr{ip, 65001, ""})
		if err != nil {
			Logger.Printf("Failed to send discovery packet: %v", err)
		}
	}

	for device := range devices {
		discoveredDevices[device.ID()] = device
	}

	return discoveredDevices, err
}

func NewDevice(conn Connection, id []byte) *Device {
	return &Device{
		Connection: conn,
		id:         id,
	}
}

func (d *Device) ID() string {
	return hex.EncodeToString(d.id)
}

func (d *Device) getset(name string, value *string) (resp string, err error) {
	err = d.Connect()
	if err != nil {
		return
	}

	if d.e == nil {
		d.e = NewEncoder(d)
	}

	if d.d == nil {
		d.d = NewDecoder(d)
	}

	tags := make(map[TagType]TagValue)
	tags[TagGetSetName] = TagValue(name)

	if value != nil {
		tags[TagGetSetValue] = TagValue(*value)
	}

	err = d.e.Encode(NewPacket(TypeGetSetReq, tags))
	if err == nil {
		var p *Packet
		p, err = d.d.Next()
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
