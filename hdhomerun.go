package hdhomerun

import (
	"io"
	"log"
	"net"
	"os"
)

var Logger = log.New(os.Stderr, "HDHomerun", log.LstdFlags)

type GetSetter interface {
	Get(name string) (string, error)
	Set(name, value string) (string, error)
}

type Connection interface {
	io.ReadWriteCloser
	Connect() error
	Location() string
}

type TcpConnection struct {
	*net.TCPConn
	ip   net.IP
	port int
}

func NewTcpConnection(ip net.IP, port int) *TcpConnection {
	return &TcpConnection{
		ip:   ip,
		port: port,
	}
}

func (conn *TcpConnection) Connect() (err error) {
	if conn.TCPConn == nil {
		conn.TCPConn, err = net.DialTCP("tcp", nil, &net.TCPAddr{conn.ip, conn.port, ""})
	}
	return
}

func (conn *TcpConnection) Location() string {
	return conn.ip.String()
}
