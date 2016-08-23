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

func deviceFromId(id string) (hdhomerun.Device, error) {
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

	for discoverResult := range hdhomerun.Discover(ip, time.Millisecond*200) {
		if discoverResult.Err != nil {
			return discoverResult.Device, discoverResult.Err
		}

		device := discoverResult.Device
		if ip != nil || discoverResult.ID.String() == id {
			return device, nil
		}
	}

	return nil, nil
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	if os.Args[1] == "discover" {
		for discoverResult := range hdhomerun.Discover(nil, time.Millisecond*200) {
			if discoverResult.Err != nil {
				fmt.Printf("Error during discovery: %v\n", discoverResult.Err)
			} else {
				fmt.Printf("hdhomerun device %s found at %s\n", discoverResult.ID, discoverResult.Addr)
			}
		}
	} else {
		device, err := deviceFromId(os.Args[1])
		if device == nil || err != nil {
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
			tuner := device.Tuner(0)
			for channel := range tuner.Scan() {
				fmt.Printf("Channel %s:%d\n", channel.Modulation, channel.Frequency)
				fmt.Printf("\tTSID: %d ONID: %d\n", channel.TSID, channel.ONID)
				for _, program := range channel.Programs {
					fmt.Printf("\t%s\n", program.Name)
				}
			}
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
