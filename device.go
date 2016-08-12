package hdhomerun

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"time"
)

type Device struct {
	id   []byte
	ip   net.IP
	port int
	conn net.Conn
	e    *Encoder
	d    *Decoder
}

func Discover(ip net.IP, timeout time.Duration) (map[string]*Device, error) {
	var wg sync.WaitGroup
	if ip == nil {
		ip = net.IPv4bcast
	}
	discovered := make(map[string]*Device)

	laddr, _ := net.ResolveUDPAddr("udp", ":0")
	conn, err := net.ListenUDP("udp", laddr)
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
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
			if discovered[raddr.String()] == nil {
				discovered[raddr.String()] = nil
				discovered[raddr.String()] = NewDevice(raddr.IP, raddr.Port, p.Tags[TagDeviceId].Value)
			}
		}
	}()

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

	wg.Wait()
	devices := make(map[string]*Device, len(discovered))
	for _, device := range discovered {
		devices[device.ID()] = device
	}
	return devices, err
}

func Connect(id, ip string, port uint16) *Device {
	d := &Device{}

	return d
}

func NewDevice(ip net.IP, port int, id []byte) *Device {
	return &Device{
		id:   id,
		ip:   ip,
		port: port,
	}
}

func (d *Device) ID() string {
	return hex.EncodeToString(d.id)
}

func (d *Device) IP() net.IP {
	return d.ip
}

func (d *Device) getset(name string, value *string) (resp string, err error) {
	if d.conn == nil {
		d.conn, err = net.DialTCP("tcp", nil, &net.TCPAddr{d.ip, 65001, ""})
		if err != nil {
			return
		}
		d.e = NewEncoder(d.conn)
		d.d = NewDecoder(d.conn)
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
			err = fmt.Errorf("Expected Get/Set reply but got %s", p.Type)
		} else {
			resp = p.Tags[TagGetSetValue].String()
			if tag, found := p.Tags[TagErrorMsg]; found {
				err = fmt.Errorf("Get/Set failed: %v", tag.String())
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
