package hdhomerun

import (
	"testing"
)

func TestParseError(t *testing.T) {
	err := ErrParseError("this is a parse error")
	if err.Error() != "this is a parse error" {
		t.Errorf("Expected 'this is a parse error' got '%s'", err.Error())
	}
}

func TestWrongPacketType(t *testing.T) {
	err := wrongPacketType(TypeDiscoverReq, TypeDiscoverRpy)
	expected := "Expected packet type Discover Request but got Discover Reply"
	if err.Error() != expected {
		t.Errorf("Expected %s got %v", expected, err)
	}
}

func TestRemoteError(t *testing.T) {
	expected := "this is a failure message"
	err := ErrRemoteError(expected)
	if err.Error() != expected {
		t.Errorf("Expected %s got %v", expected, err)
	}
}
