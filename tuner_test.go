package hdhomerun

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestParseInt(t *testing.T) {
	i, _ := parseInt("1")
	if i != 1 {
		t.Errorf("Expected 1 got %d", i)
	}
}

func TestUnmarshalStatus(t *testing.T) {
	nilType := reflect.Type(nil)
	tests := []struct {
		str            string
		expectedSignal bool
		expectedErr    reflect.Type
		expected       *TunerStatus
	}{
		{
			str:            "ch=1 lock=qam256 ss=83 snq=90 seq=100 bps=38807712 pps=1",
			expectedErr:    nilType,
			expectedSignal: true,
			expected: &TunerStatus{
				Channel:              "1",
				LockStr:              "qam256",
				Lock:                 true,
				SignalStrength:       83,
				SignalToNoiseQuality: 90,
				SymbolErrorQuality:   100,
				BitsPerSecond:        38807712,
				PacketsPerSecond:     1,
			},
		}, {
			str:            "ch=1 lock=qam256 ss=35 snq=90 seq=100 bps=38807712 pps=1",
			expectedErr:    nilType,
			expectedSignal: false,
			expected: &TunerStatus{
				Channel:              "1",
				LockStr:              "qam256", // TODO: determine an actual lock string that would start when a parenthesis
				Lock:                 true,
				SignalStrength:       35,
				SignalToNoiseQuality: 90,
				SymbolErrorQuality:   100,
				BitsPerSecond:        38807712,
				PacketsPerSecond:     1,
			},
		}, {
			str:            "ch=1 lock=(foo) ss=ff snq=90 seq=100 bps=38807712 pps=1",
			expectedErr:    reflect.TypeOf(&strconv.NumError{}),
			expectedSignal: false,
			expected: &TunerStatus{
				Channel:              "1",
				LockStr:              "(foo)", // TODO: determine an actual lock string that would start when a parenthesis
				Lock:                 false,
				SignalStrength:       35,
				SignalToNoiseQuality: 90,
				SymbolErrorQuality:   100,
				BitsPerSecond:        38807712,
				PacketsPerSecond:     1,
			},
		},
	}

	for _, test := range tests {
		received := &TunerStatus{}
		err := received.UnmarshalText([]byte(test.str))
		if reflect.TypeOf(err) != test.expectedErr {
			t.Errorf("Expected error '%v' but got '%v", test.expectedErr, reflect.TypeOf(err))
		}

		if err != nil {
			continue
		}

		if test.expected.Dump() != received.Dump() {
			t.Errorf("Expected '%s' but got '%s", test.expected.Dump(), received.Dump())
		}

		if !reflect.DeepEqual(test.expected, received) {
			t.Errorf("Expected '%s' but got '%s'", test.expected.Dump(), received.Dump())
		}

		receivedStr, _ := received.MarshalText()
		if string(receivedStr) != test.str {
			t.Errorf("Expected '%s' but got '%s'", test.str, string(receivedStr))
		}

		if received.SignalPresent() != test.expectedSignal {
			t.Errorf("Expected %v but got %v", test.expectedSignal, received.SignalPresent())
		}
	}
}

func TestStatus(t *testing.T) {
	device := newTestDevice()
	device.Send(statusRpy.p)

	tuner := device.Tuner(0)
	_, err := tuner.Status()
	if err != nil {
		t.Errorf("Did not expect an error but got %v", err)
	}
	received, _ := device.Recv()
	if !reflect.DeepEqual(statusReq.p, received) {
		t.Errorf("Expected request packet:\n%s\nGot:\n%s", statusReq.p.Dump(), received.Dump())
	}
}

func newStatusRpy(signal int, lock string) TagValue {
	return TagValue(fmt.Sprintf("ch=1 lock=%s ss=%d snq=90 seq=100 bps=38807712 pps=1", lock, signal))
}

type testGetSetter struct {
	sequence []TagValue
	err      error
}

func (t *testGetSetter) Get(name string) (TagValue, error) {
	if t.err != nil {
		return nil, t.err
	}
	resp := t.sequence[0]
	t.sequence = t.sequence[1:]
	return resp, nil
}

func (t *testGetSetter) Set(name, value string) (TagValue, error) {
	return nil, nil
}

func TestWaitForLock(t *testing.T) {
	tests := []struct {
		packetSequence []TagValue
		sendErr        error
		expectedErr    error
	}{
		{
			packetSequence: []TagValue{TagValue("foo")},
			sendErr:        wrongPacketType(TypeGetSetRpy, TypeDiscoverRpy),
			expectedErr:    wrongPacketType(TypeGetSetRpy, TypeDiscoverRpy),
		}, {
			packetSequence: []TagValue{newStatusRpy(60, "qam256")},
			expectedErr:    nil,
		}, {
			packetSequence: []TagValue{newStatusRpy(35, "qam256")},
			expectedErr:    ErrNoSignal,
		}, {
			packetSequence: []TagValue{
				newStatusRpy(60, "none"), newStatusRpy(60, "none"), newStatusRpy(60, "none"), newStatusRpy(60, "none"), newStatusRpy(60, "none"),
				newStatusRpy(60, "none"), newStatusRpy(60, "none"), newStatusRpy(60, "none"), newStatusRpy(60, "none"), newStatusRpy(60, "none"),
				newStatusRpy(60, "none"), newStatusRpy(60, "none"), newStatusRpy(60, "none"), newStatusRpy(60, "none"), newStatusRpy(60, "none"),
			},
			expectedErr: ErrTimeout,
		},
	}

	for _, test := range tests {
		tuner := newTuner(&testGetSetter{test.packetSequence, test.sendErr}, 0)
		err := tuner.WaitForLock()
		if err != test.expectedErr {
			t.Errorf("Expected err %v but got %v", test.expectedErr, err)
		}
	}
}

func TestTune(t *testing.T) {
}

func TestScan(t *testing.T) {
}
