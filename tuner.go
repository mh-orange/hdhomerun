package hdhomerun

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TunerStatus struct {
	Channel              string
	LockStr              string
	SignalPresent        bool
	LockSupported        bool
	LockUnsupported      bool
	SignalStrength       int
	SignalToNoiseQuality int
	SymbolErrorQuality   int
	BitsPerSecond        int
	PacketsPerSecond     int
}

func parseInt(str string) (int, error) {
	i, err := strconv.ParseInt(str, 10, 0)
	return int(i), err
}

func parseStatusStr(str string, s *TunerStatus) (err error) {
	for _, kv := range strings.Split(str, " ") {
		if kv != "" {
			pair := strings.Split(kv, "=")
			switch pair[0] {
			case "ch":
				s.Channel = pair[1]
			case "lock":
				s.LockStr = pair[1]
			case "ss":
				s.SignalStrength, err = parseInt(pair[1])
			case "snq":
				s.SignalToNoiseQuality, err = parseInt(pair[1])
			case "seq":
				s.SymbolErrorQuality, err = parseInt(pair[1])
			case "bps":
				s.BitsPerSecond, err = parseInt(pair[1])
			case "pps":
				s.PacketsPerSecond, err = parseInt(pair[1])
			}
		}

		if err != nil {
			break
		}
	}

	if err == nil {
		s.SignalPresent = s.SignalStrength >= 45
		if s.LockStr != "none" {
			if s.LockStr[0] == '(' {
				s.LockUnsupported = true
			} else {
				s.LockSupported = true
			}
		}
	}

	return
}

type Tuner struct {
	d *Device
	n int
}

func newTuner(d *Device, n int) *Tuner {
	return &Tuner{
		d: d,
		n: n,
	}
}

func (t *Tuner) GetTuner(name string) (string, error) {
	return t.d.Get(fmt.Sprintf("/tuner%d/%s", t.n, name))
}

func (t *Tuner) SetTuner(name, value string) (string, error) {
	return t.d.Set(fmt.Sprintf("/tuner%d/%s", t.n, name), value)
}

func (t *Tuner) Status() (status TunerStatus, err error) {
	str, err := t.GetTuner("status")
	if err == nil {
		err = parseStatusStr(str, &status)
	}
	return status, err
}

func (t *Tuner) WaitForLock() (status TunerStatus, err error) {
	time.Sleep(250 * time.Millisecond)
	timeout := time.Now().Add(2500 * time.Millisecond)

	for {
		status, err = t.Status()
		if err != nil {
			break
		}

		// TODO: this logic doesn't make sense to me, but it's how it's done in
		// the SiliconDust libhdhomerun.  Need to try and learn what the meaning
		// of this logic is... maybe contact SiliconDust about it
		if !status.SignalPresent || status.LockSupported || status.LockUnsupported {
			break
		}

		if time.Now().After(timeout) {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}

	return
}

func (t *Tuner) Tune(channel Channel) error {
	_, err := t.SetTuner("channel", fmt.Sprintf("%s:%d", channel.Modulation, channel.Frequency))
	return err
}

func (t *Tuner) Scan() chan Channel {
	visited := make(map[uint32]bool)

	ch := make(chan Channel)
	go func() {
		channelmap, err := t.GetTuner("channelmap")
		if err == nil {
			for channel := range Channels(channelmap) {
				if visited[channel.Frequency] {
					continue
				}
				visited[channel.Frequency] = true
				t.Tune(channel)
				_, err := t.WaitForLock()
				if err != nil {
					Logger.Printf("Error waiting for lock: %v", err)
				} else {
					ch <- channel
				}
			}
		} else {
			Logger.Printf("Error retrieving channel map: %v", err)
		}
		close(ch)
	}()
	return ch
}
