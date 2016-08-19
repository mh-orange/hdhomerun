package hdhomerun

import (
	"bytes"
	"encoding/hex"
	"io"
	"net"
	"sync"
	"time"
)

type Device struct {
	Connection
}

type DeviceID []byte

func (d DeviceID) String() string {
	return hex.EncodeToString(d)
}

type DiscoverResult struct {
	Device *Device
	ID     DeviceID
	Err    error
}

func Discover(ip net.IP, timeout time.Duration) chan DiscoverResult {
	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan DiscoverResult, 1)
	if ip == nil {
		ip = net.IPv4bcast
	}

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		ch <- DiscoverResult{nil, nil, err}
		return ch
	}

	go func() {
		defer close(ch)

		buffer := make([]byte, 1500)
		readBuffer := bytes.NewBuffer(make([]byte, 1500))
		readBuffer.Reset()
		decoder := NewDecoder(readBuffer)
		discovered := make(map[string]bool)

		end := time.Now().Add(timeout)
		for t := time.Now(); t.Before(end); t = time.Now() {
			conn.SetReadDeadline(time.Now().Add(timeout))

			n, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					continue
				} else if err != io.EOF {
					ch <- DiscoverResult{nil, nil, err}
				}
				break
			}

			readBuffer.Write(buffer[0:n])
			p, err := decoder.Next()

			if err != nil {
				ch <- DiscoverResult{nil, nil, err}
				break
			}

			if p.Type != TypeDiscoverRpy {
				continue
			}

			if _, found := discovered[addr.String()]; !found {
				ch <- DiscoverResult{
					Device: NewDevice(NewTCPConnection(&net.TCPAddr{addr.IP, addr.Port, addr.Zone})),
					ID:     DeviceID(p.Tags[TagDeviceId].Value),
					Err:    nil,
				}
			}
			discovered[addr.String()] = true
		}
		conn.Close()
		// Wait on the sender go routine so that we don't close the result
		// channel before it's done and cause a panic
		wg.Wait()
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
				ch <- DiscoverResult{nil, nil, err}
			}
		}
		wg.Done()
	}()

	return ch
}

func Connect(ip net.IP, port int) (*Device, error) {
	device := NewDevice(NewTCPConnection(&net.TCPAddr{ip, port, ""}))
	return device, device.Connect()
}

func NewDevice(conn Connection) *Device {
	return &Device{
		Connection: conn,
	}
}

func (d *Device) getset(name string, value *string) (resp TagValue, err error) {
	err = d.Connect()
	tags := make(map[TagType]TagValue)
	tags[TagGetSetName] = TagValue(name)

	if value != nil {
		tags[TagGetSetValue] = TagValue(*value)
	}

	if err == nil {
		err = d.Send(NewPacket(TypeGetSetReq, tags))
	}

	if err == nil {
		var p *Packet
		p, err = d.Recv()
		if p.Type != TypeGetSetRpy {
			err = wrongPacketType(TypeGetSetRpy, p.Type)
		} else {
			resp = p.Tags[TagGetSetValue].Value
			if tag, found := p.Tags[TagErrorMsg]; found {
				err = ErrRemoteError(tag.String())
			}
		}
	}
	return
}

func (d *Device) Get(name string) (TagValue, error) {
	return d.getset(name, nil)
}

func (d *Device) Set(name, value string) (TagValue, error) {
	return d.getset(name, &value)
}

func (d *Device) Tuner(n int) *Tuner {
	return newTuner(d, n)
}

func (d *Device) Connect() error {
	if conn, ok := d.Connection.(Connectable); ok {
		return conn.Connect()
	}
	return nil
}

func (d *Device) Close() error {
	if conn, ok := d.Connection.(Closeable); ok {
		return conn.Close()
	}
	return nil
}

func (d *Device) Addr() net.Addr {
	if conn, ok := d.Connection.(Addressable); ok {
		return conn.RemoteAddr()
	}
	return nil
}
