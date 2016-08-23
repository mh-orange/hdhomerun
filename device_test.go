package hdhomerun

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func newTestDevice() *GenericDevice {
	var buffer bytes.Buffer
	return NewGenericDevice(&buffer)
}

func TestGetSet(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		reply         *Packet
		expectedValue TagValue
		expectedErr   reflect.Type
	}{
		{
			name:          "help",
			reply:         getRpy.p,
			expectedValue: getRpy.p.Tags[TagGetSetValue].Value,
		}, {
			name:          "/tuner0/channel",
			value:         "auto:849000000",
			reply:         setRpy.p,
			expectedValue: setRpy.p.Tags[TagGetSetValue].Value,
		}, {
			name:          "help",
			reply:         discoverRpy.p,
			expectedValue: setRpy.p.Tags[TagGetSetValue].Value,
			expectedErr:   reflect.TypeOf(ErrWrongPacketType("")),
		}, {
			name:        "help",
			reply:       getRpyErr.p,
			expectedErr: reflect.TypeOf(ErrRemoteError("")),
		},
	}

	for _, test := range tests {
		d := newTestDevice()
		d.encoder.Encode(test.reply)

		var value TagValue
		var err error
		if test.value == "" {
			value, err = d.Get(test.name)
		} else {
			value, err = d.Set(test.name, test.value)
		}

		if reflect.TypeOf(err) != test.expectedErr {
			t.Errorf("Expected error %v but got %v", test.expectedErr, reflect.TypeOf(err))
		}

		if err != nil {
			continue
		}

		if !reflect.DeepEqual(value, test.expectedValue) {
			t.Errorf("Expected return value of %s but got %s", test.expectedValue, value)
		}
	}
}

func TestDiscover(t *testing.T) {
	tests := []struct {
		reply   testPacket
		devices []string
		err     reflect.Type
	}{
		{
			reply:   discoverRpy,
			devices: []string{hex.EncodeToString(discoverRpy.p.Tags[TagDeviceId].Value)},
		}, {
			reply:   discoverReq,
			devices: []string{},
		}, {
			reply: testPacket{
				p: nil,
				b: []byte{
					0x00,
				},
			},
			devices: []string{},
			err:     reflect.TypeOf(fmt.Errorf("")),
		},
	}

	for _, test := range tests {
		listener, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 65001})
		go func() {
			listener.SetReadDeadline(time.Now().Add(time.Second))
			_, addr, _ := listener.ReadFromUDP(make([]byte, 1024))
			listener.WriteTo(test.reply.b, addr)
			listener.Close()
		}()

		devices := make([]string, 0)
		for result := range Discover(net.IP{127, 0, 0, 1}, time.Second) {
			if reflect.TypeOf(result.Err) != test.err {
				t.Errorf("Expected error type %v but got %v(%v)", test.err, reflect.TypeOf(result.Err), result.Err)
			}

			if result.Err != nil {
				continue
			}
			devices = append(devices, result.ID.String())
		}

		if !reflect.DeepEqual(devices, test.devices) {
			t.Errorf("Expected devices %v but got %v", test.devices, devices)
		}
	}
}

func TestTCPDevice(t *testing.T) {
	tests := []struct {
		txPackets []*Packet
		rxPackets []*Packet
	}{
		{
			txPackets: []*Packet{getReq.p},
			rxPackets: []*Packet{getRpy.p},
		},
	}

	for _, test := range tests {
		var wg sync.WaitGroup
		listener, _ := net.ListenTCP("tcp", &net.TCPAddr{net.IP{127, 0, 0, 1}, 65001, ""})

		wg.Add(1)
		go func() {
			conn, _ := listener.Accept()
			device := NewGenericDevice(conn)
			for i, expectedTx := range test.txPackets {
				receivedTx, _ := device.decoder.Next()
				if !reflect.DeepEqual(expectedTx, receivedTx) {
					t.Errorf("Expected:\n%s\nGot:\n%s\n", expectedTx.Dump(), receivedTx.Dump())
				}
				device.encoder.Encode(test.rxPackets[i])
			}
			conn.Close()
			wg.Done()
		}()

		d, _ := ConnectTCP(&net.TCPAddr{net.IP{127, 0, 0, 1}, 65001, ""})
		for i, expectedRx := range test.rxPackets {
			d.encoder.Encode(test.txPackets[i])
			receivedRx, _ := d.GenericDevice.decoder.Next()
			if !reflect.DeepEqual(expectedRx, receivedRx) {
				t.Errorf("Expected:\n%s\nGot:\n%s\n", expectedRx.Dump(), receivedRx.Dump())
			}
		}
		d.Close()
		wg.Wait()
	}
}
