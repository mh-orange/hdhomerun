package hdhomerun

import (
	"fmt"
)

type ErrWrongPacketType string

func wrongPacketType(expected, received PacketType) error {
	return ErrWrongPacketType(fmt.Sprintf("Expected %s reply but got %s", expected, received))
}

func (e ErrWrongPacketType) Error() string {
	return string(e)
}

type ErrRemoteError string

func (e ErrRemoteError) Error() string {
	return string(e)
}
