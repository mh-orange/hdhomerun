package hdhomerun

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

type testConnection struct {
	reader     *bytes.Buffer
	writer     *bytes.Buffer
	connectErr error
}

func (t testConnection) Read(p []byte) (int, error) {
	return t.reader.Read(p)
}

func (t testConnection) Write(p []byte) (int, error) {
	return t.writer.Write(p)
}

func (t testConnection) Close() error {
	return nil
}

func (t testConnection) Connect() error {
	return t.connectErr
}

func (t testConnection) Location() string {
	return ""
}

func mockConnection() testConnection {
	return testConnection{
		reader: &bytes.Buffer{},
		writer: &bytes.Buffer{},
	}
}

func TestID(t *testing.T) {
	mock := mockConnection()
	d := NewDevice(mock, []byte{0x01, 0x02, 0x03, 0x04})
	if d.ID() != "01020304" {
		t.Errorf("Expected ID to be %s but got %s", "01020304", d.ID())
	}
}

func TestConnectErr(t *testing.T) {
	mock := mockConnection()
	errMsg := "This is a test error"
	mock.connectErr = fmt.Errorf(errMsg)
	d := NewDevice(mock, []byte{0x01, 0x02, 0x03, 0x04})
	_, err := d.Get("help")
	if err == nil || err.Error() != errMsg {
		t.Errorf("Expected connect error %s but got %v", errMsg, err)
	}
}

func TestWrongPacketType(t *testing.T) {
	mock := mockConnection()
	mock.reader.Write(discoverRpy.b)
	d := NewDevice(mock, []byte{0x01, 0x02, 0x03, 0x04})
	_, err := d.Get("help")
	if _, ok := err.(ErrWrongPacketType); !ok {
		t.Errorf("Expected error ErrWrongPacketType but got %T", err)
	}
}

func TestRemoteError(t *testing.T) {
	mock := mockConnection()
	mock.reader.Write(getRpyErr.b)
	d := NewDevice(mock, []byte{0x01, 0x02, 0x03, 0x04})
	_, err := d.Get("help")
	if _, ok := err.(ErrRemoteError); !ok {
		t.Errorf("Expected error ErrREmoteError but got %T", err)
	}
}

func TestGet(t *testing.T) {
	mock := mockConnection()
	mock.reader.Write(getRpy.b)
	d := NewDevice(mock, []byte{0x01, 0x02, 0x03, 0x04})
	value, _ := d.Get("help")
	if !reflect.DeepEqual(mock.writer.Bytes(), getReq.b) {
		t.Errorf("Device get failed\nExpected:\n%s\nReceived:\n%s\n", hex.Dump(mock.reader.Bytes()), getReq.b)
	}

	if value != getRpy.p.Tags[TagGetSetValue].String() {
		t.Errorf("Device get failed\nExpected: %s Received: %s", getRpy.p.Tags[TagGetSetValue], value)
	}
}

func TestSet(t *testing.T) {
	mock := mockConnection()
	mock.reader.Write(setRpy.b)
	d := NewDevice(mock, []byte{0x01, 0x02, 0x03, 0x04})
	d.Set("/tuner0/channel", "auto:849000000")
	if !reflect.DeepEqual(mock.writer.Bytes(), setReq.b) {
		t.Errorf("Device set failed\nExpected:\n%s\nReceived:\n%s\n", hex.Dump(setReq.b), hex.Dump(mock.writer.Bytes()))
	}
}

func TestTuner(t *testing.T) {
	// this is almost ridiculous
	mock := mockConnection()
	d := NewDevice(mock, []byte{0x01, 0x02, 0x03, 0x04})
	tuner := d.Tuner(1)
	if tuner.n != 1 {
		t.Errorf("Expected tuner number to be 1 but got %d", tuner.n)
	}

	if tuner.d != d {
		t.Errorf("Expected device to be %v but got %v", d, tuner.d)
	}

}
