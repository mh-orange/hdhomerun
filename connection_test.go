package hdhomerun

import (
	"net"
	"reflect"
	"sync"
	"testing"
)

func TestTCPConnection(t *testing.T) {
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
			io := NewIOConnection(conn)
			for i, expectedTx := range test.txPackets {
				receivedTx, _ := io.Recv()
				if !reflect.DeepEqual(expectedTx, receivedTx) {
					t.Errorf("Expected:\n%s\nGot:\n%s\n", expectedTx.Dump(), receivedTx.Dump())
				}
				io.Send(test.rxPackets[i])
			}
			conn.Close()
			wg.Done()
		}()

		d := NewDevice(NewTCPConnection(&net.TCPAddr{net.IP{127, 0, 0, 1}, 65001, ""}))
		d.Connect()
		for i, expectedRx := range test.rxPackets {
			d.Send(test.txPackets[i])
			receivedRx, _ := d.Recv()
			if !reflect.DeepEqual(expectedRx, receivedRx) {
				t.Errorf("Expected:\n%s\nGot:\n%s\n", expectedRx.Dump(), receivedRx.Dump())
			}
		}
		d.Close()
		wg.Wait()
	}
}
