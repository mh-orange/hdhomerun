package hdhomerun

import (
	"fmt"
	"net"
)

type Device struct {
	c   net.Conn
	err error
	e   *Encoder
	d   *Decoder
}

func Connect(id, ip string, port uint16) *Device {
	d := &Device{}

	d.c, d.err = net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	d.e = NewEncoder(d.c)
	d.d = NewDecoder(d.c)

	return d
}

func (d *Device) Get(name string) {
	// send packet
	// receive response
}

func (d *Device) Set(name, value string) {
	// send packet
	// receive response
}
