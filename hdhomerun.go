package hdhomerun

import (
	"log"
	"os"
)

var Logger = log.New(os.Stderr, "HDHomerun", log.LstdFlags)
