package hdhomerun

import (
	"fmt"
	"net"
)

type Device struct {
	c   net.Conn
	err error
	e   *encoder
	d   *decoder
}

func Connect(id, ip string, port uint16) *Device {
	d := &Device{}

	d.c, d.err = net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	d.e = newEncoder(d.c)
	d.d = newDecoder(d.c)

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
