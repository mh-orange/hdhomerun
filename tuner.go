package hdhomerun

import (
	"fmt"
	"strings"
	"time"
)

type TunerStatus struct {
	Channel              string
	LockStr              string
	LockSupported        bool
	SignalStrength       int
	SignalToNoiseQuality int
	SymbolErrorQuality   int
	BitsPerSecond        int
	PacketsPerSecond     int
}

func (ts *TunerStatus) MarshalText() ([]byte, error) {
	return []byte(ts.Dump()), nil
}

func (ts *TunerStatus) Dump() string {
	return fmt.Sprintf("ch=%s lock=%s ss=%d snq=%d seq=%d bps=%d pps=%d",
		ts.Channel,
		ts.LockStr,
		ts.SignalStrength,
		ts.SignalToNoiseQuality,
		ts.SymbolErrorQuality,
		ts.BitsPerSecond,
		ts.PacketsPerSecond,
	)
}

func (ts *TunerStatus) UnmarshalText(text []byte) (err error) {
	str := string(text)

	for _, kv := range strings.Split(str, " ") {
		if kv != "" {
			pair := strings.Split(kv, "=")
			switch pair[0] {
			case "ch":
				ts.Channel = pair[1]
			case "lock":
				ts.LockStr = pair[1]
			case "ss":
				ts.SignalStrength, err = parseInt(pair[1])
			case "snq":
				ts.SignalToNoiseQuality, err = parseInt(pair[1])
			case "seq":
				ts.SymbolErrorQuality, err = parseInt(pair[1])
			case "bps":
				ts.BitsPerSecond, err = parseInt(pair[1])
			case "pps":
				ts.PacketsPerSecond, err = parseInt(pair[1])
			}
		}

		if err != nil {
			break
		}
	}

	if err == nil {
		if ts.LockStr != "none" {
			ts.LockSupported = true
		}
	}

	return
}

func (ts *TunerStatus) SignalPresent() bool {
	return ts.SignalStrength >= 45
}

type Tuner struct {
	d GetSetter
	n int
}

func newTuner(d GetSetter, n int) *Tuner {
	return &Tuner{
		d: d,
		n: n,
	}
}

func (t *Tuner) GetTuner(name string) (TagValue, error) {
	return t.d.Get(fmt.Sprintf("/tuner%d/%s", t.n, name))
}

func (t *Tuner) SetTuner(name, value string) (TagValue, error) {
	return t.d.Set(fmt.Sprintf("/tuner%d/%s", t.n, name), value)
}

func (t *Tuner) Status() (*TunerStatus, error) {
	status := &TunerStatus{}

	value, err := t.GetTuner("status")
	if err == nil {
		err = status.UnmarshalText(value)
	}
	return status, err
}

func (t *Tuner) WaitForLock() (err error) {
	var status *TunerStatus

	time.Sleep(250 * time.Millisecond)
	timeout := time.Now().Add(2500 * time.Millisecond)

	for {
		status, err = t.Status()
		if err != nil {
			break
		}

		if !status.SignalPresent() {
			err = ErrNoSignal
			break
		} else if status.LockSupported {
			err = ErrLockNotSupported
			break
		}

		if time.Now().After(timeout) {
			err = ErrTimeout
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

func (t *Tuner) StreamInfo() (TagValue, error) {
	return t.GetTuner("streaminfo")
}

func (t *Tuner) Scan() chan Channel {
	visited := make(map[uint32]bool)

	ch := make(chan Channel)
	go func() {
		channelmap, err := t.GetTuner("channelmap")
		if err == nil {
			for channel := range Channels(channelmap.String()) {
				if visited[channel.Frequency] {
					continue
				}
				visited[channel.Frequency] = true
				t.Tune(channel)
				err := t.WaitForLock()
				if err != nil {
					if err != ErrLockNotSupported && err != ErrNoSignal {
						Logger.Printf("Error waiting for lock: %v", err)
					}
					continue
				}

				timeout := time.Now().Add(5 * time.Second)
				var status *TunerStatus
				for {
					status, err = t.Status()
					if err != nil {
						break
					}

					if status.SymbolErrorQuality == 100 || time.Now().After(timeout) {
						break
					}

					time.Sleep(250 * time.Millisecond)
				}

				if err == nil {
					channel.Name = fmt.Sprintf("%s:%d", status.LockStr, channel.Number)
					var si TagValue
					si, err = t.StreamInfo()
					if err == nil {
						channel.UnmarshalText(si)
					}
				}

				if err == nil {
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
