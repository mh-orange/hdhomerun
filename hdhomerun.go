package hdhomerun

import (
	"log"
	"os"
	"strconv"
)

var Logger = log.New(os.Stderr, "HDHomerun", log.LstdFlags)

func parseInt(str string) (int, error) {
	i, err := strconv.ParseInt(str, 10, 0)
	return int(i), err
}
