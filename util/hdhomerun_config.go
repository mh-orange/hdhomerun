package main

import (
	"fmt"
	"github.com/abates/hdhomerun"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

func usage() {
	programName := path.Base(os.Args[0])
	fmt.Printf("Usage:\n")
	fmt.Printf("\t%s discover\n", programName)
	fmt.Printf("\t%s <id> get help\n", programName)
	fmt.Printf("\t%s <id> get <item>\n", programName)
	fmt.Printf("\t%s <id> set <item> <value>\n", programName)
	fmt.Printf("\t%s <id> scan <tuner> [<filename>]\n", programName)
	fmt.Printf("\t%s <id> save <tuner> <filename>\n", programName)
	fmt.Printf("\t%s <id> upgrade <filename>\n", programName)
	os.Exit(0)
}

func deviceFromId(id string) (*hdhomerun.Device, error) {
	var ip net.IP
	if strings.Contains(id, ".") {
		if strings.Contains(id, ":") {
			return nil, fmt.Errorf("Multicast not yet implemented")
		} else if strings.Contains(id, "-") {
			return nil, fmt.Errorf("Tuner number not yet implemented")
		} else {
			ipAddr, err := net.ResolveIPAddr("ip4", id)
			if err != nil {
				return nil, err
			}
			ip = ipAddr.IP
		}
	}
	devices, err := hdhomerun.Discover(ip, time.Millisecond*200)
	if ip != nil {
		var device *hdhomerun.Device
		for _, device = range devices {
			return device, err
		}

		if err == nil && device == nil {
			err = fmt.Errorf("Timeout connecting to to %s", id)
		}
	}
	return devices[id], err
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	if os.Args[1] == "discover" {
		devices, err := hdhomerun.Discover(nil, time.Millisecond*200)
		if err != nil {
			fmt.Printf("Discovery failed: %v\n", err)
			os.Exit(1)
		}

		for _, device := range devices {
			fmt.Printf("hdhomerun device %s found at %v\n", device.ID(), device.IP())
		}
	} else {
		device, err := deviceFromId(os.Args[1])
		if err != nil {
			fmt.Printf("Could not connect to device: %v\n", err)
			os.Exit(1)
		}

		switch os.Args[2] {
		case "get":
			if len(os.Args) < 4 {
				usage()
			}
			resp, err := device.Get(os.Args[3])
			if err != nil {
				fmt.Printf("Failed to get value: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", resp)
		case "set":
			if len(os.Args) < 5 {
				usage()
			}
			resp, err := device.Set(os.Args[3], os.Args[4])
			if err != nil {
				fmt.Printf("Failed to set value: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", resp)
		case "scan":
			fmt.Printf("Not implemented\n")
		case "save":
			fmt.Printf("Not implemented\n")
		case "upgrade":
			fmt.Printf("Not implemented\n")
		default:
			fmt.Printf("Unknown command %s\n", os.Args[3])
			os.Exit(1)
		}
	}
}
