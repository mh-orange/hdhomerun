package hdhomerun

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type testAddr string

func (t testAddr) Network() string {
	return string(t)
}

func (t testAddr) String() string {
	return string(t)
}

func newTestDevice() (*Device, *mockConnection) {
	mock := newMockConnection()
	return NewDevice(mock, []byte{0x01, 0x02, 0x03, 0x04}, "test address"), mock
}

func TestID(t *testing.T) {
	d, _ := newTestDevice()
	if d.ID() != "01020304" {
		t.Errorf("Expected ID to be %s but got %s", "01020304", d.ID())
	}
}

func TestReceivedWrongPacketType(t *testing.T) {
	d, mock := newTestDevice()
	mock.reader.Encode(discoverRpy.p)
	_, err := d.Get("help")
	if _, ok := err.(ErrWrongPacketType); !ok {
		t.Errorf("Expected error ErrWrongPacketType but got %T", err)
	}
}

func TestReceivedRemoteError(t *testing.T) {
	d, mock := newTestDevice()
	mock.reader.Encode(getRpyErr.p)
	_, err := d.Get("help")
	if _, ok := err.(ErrRemoteError); !ok {
		t.Errorf("Expected error ErrREmoteError but got %T", err)
	}
}

func TestGet(t *testing.T) {
	d, mock := newTestDevice()
	mock.reader.Encode(getRpy.p)
	value, _ := d.Get("help")
	if !bytes.Equal(value, getRpy.p.Tags[TagGetSetValue].Value) {
		t.Errorf("Device get failed.\nExpected:\n%s\nReceived:\n%s\n", getRpy.p.Tags[TagGetSetValue].String(), value)
	}
}

func TestSet(t *testing.T) {
	d, mock := newTestDevice()
	mock.reader.Encode(setRpy.p)
	d.Set("/tuner0/channel", "auto:849000000")
	expected := NewPacket(TypeGetSetReq, map[TagType]TagValue{TagGetSetName: TagValue("/tuner0/channel"), TagGetSetValue: TagValue("auto:849000000")})
	received, _ := mock.writer.Next()
	if !reflect.DeepEqual(received, expected) {
		t.Errorf("Expected packet to be:\n%s\nBut got:\n%s\n", expected.Dump(), received.Dump())
	}
}

func TestTuner(t *testing.T) {
	// this is almost ridiculous
	d, _ := newTestDevice()
	tuner := d.Tuner(1)
	if tuner.n != 1 {
		t.Errorf("Expected tuner number to be 1 but got %d", tuner.n)
	}

	if tuner.d != d {
		t.Errorf("Expected device to be %v but got %v", d, tuner.d)
	}
}

func TestDiscovery(t *testing.T) {
	conn := newMockDiscoverConn()
	listenUDP = func() (discoverConn, error) {
		return conn, nil
	}

	conn.reader.Encode(discoverRpy.p)
	count := 0
	devices, _ := Discover(nil, time.Millisecond*100)
	for range devices {
		count++
	}

	if count != 1 {
		t.Errorf("Expected to discover one device but got %d", count)
	}
}

func TestDiscoveryWrongPacket(t *testing.T) {
	conn := newMockDiscoverConn()
	listenUDP = func() (discoverConn, error) {
		return conn, nil
	}

	conn.reader.Encode(discoverReq.p)
	count := 0
	devices, _ := Discover(nil, time.Millisecond*100)
	for range devices {
		count++
	}

	if count != 0 {
		t.Errorf("Expected to discover no devices but got %d", count)
	}
}

type errTimeout string

func (e errTimeout) Timeout() bool {
	return true
}

func (e errTimeout) Temporary() bool {
	return false
}

func (e errTimeout) Error() string {
	return string(e)
}

func TestDiscoveryConnecionError(t *testing.T) {
	expectedErr := fmt.Errorf("connection error")
	listenUDP = func() (discoverConn, error) {
		return nil, expectedErr
	}

	_, err := Discover(nil, time.Millisecond*100)

	if err != expectedErr {
		t.Errorf("Expected error to be %v but got %v", expectedErr, err)
	}
}

func TestDiscoveryTimeout(t *testing.T) {
	conn := newMockDiscoverConn()
	conn.readErr = errTimeout("timeout")

	listenUDP = func() (discoverConn, error) {
		return conn, nil
	}
	count := 0

	devices, _ := Discover(nil, time.Millisecond*100)
	for range devices {
		count++
	}

	if count != 0 {
		t.Errorf("Expected to discover no devices but got %d", count)
	}
}
