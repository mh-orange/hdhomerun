package hdhomerun

import (
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

func TestParseStatusStr(t *testing.T) {
	nilType := reflect.Type(nil)
	tests := []struct {
		str         string
		expectedErr reflect.Type
		expected    *TunerStatus
	}{
		{
			str:         "ch=1 lock=qam256 ss=83 snq=90 seq=100 bps=38807712 pps=1",
			expectedErr: nilType,
			expected: &TunerStatus{
				Channel:              "1",
				LockStr:              "qam256",
				SignalPresent:        true,
				LockSupported:        true,
				LockUnsupported:      false,
				SignalStrength:       83,
				SignalToNoiseQuality: 90,
				SymbolErrorQuality:   100,
				BitsPerSecond:        38807712,
				PacketsPerSecond:     1,
			},
		}, {
			str:         "ch=1 lock=qam256 ss=35 snq=90 seq=100 bps=38807712 pps=1",
			expectedErr: nilType,
			expected: &TunerStatus{
				Channel:              "1",
				LockStr:              "qam256", // TODO: determine an actual lock string that would start when a parenthesis
				SignalPresent:        false,
				LockSupported:        true,
				LockUnsupported:      false,
				SignalStrength:       35,
				SignalToNoiseQuality: 90,
				SymbolErrorQuality:   100,
				BitsPerSecond:        38807712,
				PacketsPerSecond:     1,
			},
		}, {
			str:         "ch=1 lock=(foo) ss=35 snq=90 seq=100 bps=38807712 pps=1",
			expectedErr: nilType,
			expected: &TunerStatus{
				Channel:              "1",
				LockStr:              "(foo)", // TODO: determine an actual lock string that would start when a parenthesis
				SignalPresent:        false,
				LockSupported:        false,
				LockUnsupported:      true,
				SignalStrength:       35,
				SignalToNoiseQuality: 90,
				SymbolErrorQuality:   100,
				BitsPerSecond:        38807712,
				PacketsPerSecond:     1,
			},
		}, {
			str:         "ch=1 lock=(foo) ss=ff snq=90 seq=100 bps=38807712 pps=1",
			expectedErr: reflect.TypeOf(&strconv.NumError{}),
			expected: &TunerStatus{
				Channel:              "1",
				LockStr:              "(foo)", // TODO: determine an actual lock string that would start when a parenthesis
				SignalPresent:        false,
				LockSupported:        false,
				LockUnsupported:      true,
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
		err := parseStatusStr(test.str, received)
		if reflect.TypeOf(err) != test.expectedErr {
			t.Errorf("Expected error '%v' but got '%v", test.expectedErr, reflect.TypeOf(err))
		}

		if test.expectedErr == nil {
			if test.expected.Dump() != received.Dump() {
				t.Errorf("Expected '%s' but got '%s", test.expected.Dump(), received.Dump())
			}

			if !reflect.DeepEqual(test.expected, received) {
				t.Errorf("Expected '%s' but got '%s", test.expected.Dump(), received.Dump())
			}
		}
	}
}

func TestStatus(t *testing.T) {
}

func TestWaitForLock(t *testing.T) {
}

func TestTune(t *testing.T) {
}

func TestScan(t *testing.T) {
}