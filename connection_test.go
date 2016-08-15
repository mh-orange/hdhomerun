package hdhomerun

import (
	"bytes"
	"net"
	"reflect"
	"testing"
	"time"
)

type mockStreamConnection struct {
	reader *bytes.Buffer
	writer *bytes.Buffer
}

func (s *mockStreamConnection) Read(b []byte) (int, error) {
	return s.reader.Read(b)
}

func (s *mockStreamConnection) Write(b []byte) (int, error) {
	return s.writer.Write(b)
}

func newMockStreamConnection() *mockStreamConnection {
	return &mockStreamConnection{
		reader: &bytes.Buffer{},
		writer: &bytes.Buffer{},
	}
}

type mockConnection struct {
	*IOConnection
	*mockStreamConnection
	reader  *Encoder
	writer  *Decoder
	readErr error
}

func newMockConnection() *mockConnection {
	rw := newMockStreamConnection()

	return &mockConnection{
		IOConnection:         NewIOConnection(rw),
		mockStreamConnection: rw,
		reader:               NewEncoder(rw.reader),
		writer:               NewDecoder(rw.writer),
	}
}

func (conn *mockConnection) Connect() error {
	return nil
}

func (conn *mockConnection) Close() error {
	return nil
}

type mockDiscoverConn struct {
	*mockConnection
}

func newMockDiscoverConn() *mockDiscoverConn {
	return &mockDiscoverConn{newMockConnection()}
}

func (conn *mockDiscoverConn) ReadFrom(b []byte) (int, net.Addr, error) {
	if conn.readErr != nil {
		return 0, nil, conn.readErr
	}
	n, err := conn.mockConnection.Read(b)
	return n, testAddr("test address"), err
}

func (conn *mockDiscoverConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	return conn.mockConnection.Write(b)
}

func (conn *mockDiscoverConn) SetReadDeadline(time.Time) error {
	return nil
}

func TestIOSend(t *testing.T) {
	conn := newMockConnection()
	conn.Send(discoverReq.p)
	received, _ := conn.writer.Next()
	if !reflect.DeepEqual(discoverReq.p, received) {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", discoverReq.p.Dump(), received.Dump())
	}
}

/*func TestTcpConnectionsSend(t *testing.T) {
	rwc := newTestReadWriteCloser()
	dialTCP = func(net string, laddr, raddr *net.TCPAddr) (io.ReadWriteCloser, error) {
		return rwc, nil
	}

	conn := NewTcpConnection(nil, 0)
	conn.Send(getReq.p)

	received := rwc.writer.Bytes()
	if !bytes.Equal(received, getReq.b) {
		t.Errorf("Expected:\n%v\nGot:\n%v\n", hex.Dump(getReq.b), hex.Dump(received))
	}
}*/
