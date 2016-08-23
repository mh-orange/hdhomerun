package hdhomerun

import (
	"bytes"
	"encoding/hex"
	"io"
	"net"
	"sync"
	"time"
)

type DeviceID []byte

func (d DeviceID) String() string {
	return hex.EncodeToString(d)
}

type Device interface {
	Get(string) (TagValue, error)
	Set(string, string) (TagValue, error)
	Tuner(int) *Tuner
}

type DiscoverResult struct {
	Device Device
	ID     DeviceID
	Addr   net.Addr
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
		ch <- DiscoverResult{nil, nil, nil, err}
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
					ch <- DiscoverResult{nil, nil, nil, err}
				}
				break
			}

			readBuffer.Write(buffer[0:n])
			p, err := decoder.Next()

			if err != nil {
				ch <- DiscoverResult{nil, nil, nil, err}
				break
			}

			if p.Type != TypeDiscoverRpy {
				continue
			}

			if _, found := discovered[addr.String()]; !found {
				ch <- DiscoverResult{
					Device: NewTCPDevice(&net.TCPAddr{addr.IP, addr.Port, addr.Zone}),
					Addr:   addr,
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
				ch <- DiscoverResult{nil, nil, nil, err}
			}
		}
		wg.Done()
	}()

	return ch
}

type TCPDevice struct {
	*net.TCPConn
	*GenericDevice
	addr *net.TCPAddr
}

func ConnectTCP(addr *net.TCPAddr) (d *TCPDevice, err error) {
	d = NewTCPDevice(addr)
	return d, d.connect()
}

func NewTCPDevice(addr *net.TCPAddr) *TCPDevice {
	d := &TCPDevice{
		addr: addr,
	}
	d.GenericDevice = NewGenericDevice(d)
	return d
}

func (d *TCPDevice) connect() (err error) {
	d.TCPConn, err = net.DialTCP("tcp", nil, d.addr)
	return err
}

type GenericDevice struct {
	encoder *Encoder
	decoder *Decoder
}

func NewGenericDevice(rw io.ReadWriter) *GenericDevice {
	return &GenericDevice{
		encoder: NewEncoder(rw),
		decoder: NewDecoder(rw),
	}
}

func (d *GenericDevice) getset(name string, value *string) (resp TagValue, err error) {
	tags := make(map[TagType]TagValue)
	tags[TagGetSetName] = TagValue(name)

	if value != nil {
		tags[TagGetSetValue] = TagValue(*value)
	}

	if err == nil {
		err = d.encoder.Encode(NewPacket(TypeGetSetReq, tags))
	}

	if err == nil {
		var p *Packet
		p, err = d.decoder.Next()
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

func (d *GenericDevice) Get(name string) (TagValue, error) {
	return d.getset(name, nil)
}

func (d *GenericDevice) Set(name, value string) (TagValue, error) {
	return d.getset(name, &value)
}

func (d *GenericDevice) Tuner(n int) *Tuner {
	return newTuner(d, n)
}
