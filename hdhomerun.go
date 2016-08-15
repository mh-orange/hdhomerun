package hdhomerun

import (
	"log"
	"os"
)

var Logger = log.New(os.Stderr, "HDHomerun", log.LstdFlags)

type GetSetter interface {
	Get(name string) (string, error)
	Set(name, value string) (string, error)
}
