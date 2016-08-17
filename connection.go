package hdhomerun

import (
	"io"
	"net"
)

type Connection interface {
	Send(*Packet) error
	Recv() (*Packet, error)
	RemoteAddr() net.Addr
}

type Connectable interface {
	Connect() error
}

type Closeable interface {
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

func (conn *IOConnection) RemoteAddr() net.Addr {
	return nil
}

type TCPConnection struct {
	*net.TCPConn
	*IOConnection
	addr *net.TCPAddr
}

func NewTCPConnection(addr *net.TCPAddr) *TCPConnection {
	return &TCPConnection{
		addr: addr,
	}
}

func (conn *TCPConnection) Connect() (err error) {
	conn.TCPConn, err = net.DialTCP("tcp", nil, conn.addr)
	if err == nil {
		conn.IOConnection = NewIOConnection(conn.TCPConn)
	}
	return err
}
