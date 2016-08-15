package hdhomerun

import (
	"io"
	"net"
)

type Connection interface {
	Connect() error
	Send(p *Packet) error
	Recv() (*Packet, error)
	Close() error
}

type IOConnection struct {
	encoder *Encoder
	decoder *Decoder
}

func NewIOConnection(rw io.ReadWriter) *IOConnection {
	return &IOConnection{
		encoder: NewEncoder(rw),
		decoder: NewDecoder(rw),
	}
}

func (conn *IOConnection) Send(p *Packet) error {
	return conn.encoder.Encode(p)
}

func (conn *IOConnection) Recv() (p *Packet, err error) {
	return conn.decoder.Next()
}

type TCPConnection struct {
	*net.TCPConn
	*IOConnection
	address string
}

func NewTCPConnection(address string) Connection {
	conn := &TCPConnection{
		address: address,
	}

	conn.IOConnection = NewIOConnection(conn)
	return conn
}

func (conn *TCPConnection) Connect() (err error) {
	if conn.TCPConn == nil {
		var addr *net.TCPAddr
		addr, err = net.ResolveTCPAddr("tcp", conn.address)
		if err == nil {
			conn.TCPConn, err = net.DialTCP("tcp", nil, addr)
		}
	}
	return
}
