package hdhomerun

import (
	"fmt"
)

var (
	ErrTimeout  = fmt.Errorf("Timeout")
	ErrNoSignal = fmt.Errorf("Signal Strength Too Low")
)

type ErrCommunicationError struct {
	ConnectionError error
}

type ErrParseError string

func (e ErrParseError) Error() string {
	return string(e)
}

type ErrWrongPacketType string

func wrongPacketType(expected, received PacketType) error {
	return ErrWrongPacketType(fmt.Sprintf("Expected packet type %s but got %s", expected, received))
}

func (e ErrWrongPacketType) Error() string {
	return string(e)
}

type ErrRemoteError string

func (e ErrRemoteError) Error() string {
	return string(e)
}
