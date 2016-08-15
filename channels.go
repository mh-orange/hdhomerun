package hdhomerun

import (
	"fmt"
	"sync"
)

type Channel struct {
	Number     int
	Frequency  uint32
	Modulation string
	Name       string
}

type channelRange struct {
	Start     int
	End       int
	Frequency uint32
	Spacing   uint32
}

type channelMap []channelRange

func channelFrequencyRound(frequency, resolution uint32) uint32 {
	frequency += resolution / 2
	return (frequency / resolution) * resolution
}

func channelFrequencyRoundNormal(frequency uint32) uint32 {
	return channelFrequencyRound(frequency, 125000)
}

func (cr channelRange) Channels() []Channel {
	ch := make(chan Channel)
	channels := make([]Channel, 0)
	go func() { cr.channels(ch); close(ch) }()
	for c := range ch {
		channels = append(channels, c)
	}
	return channels
}

func (cr channelRange) channels(txCh chan Channel) {
	for i := cr.Start; i <= cr.End; i++ {
		txCh <- Channel{
			Number:     i,
			Frequency:  channelFrequencyRoundNormal(cr.Frequency + (uint32(i-cr.Start) * cr.Spacing)),
			Modulation: "auto",
			Name:       fmt.Sprintf("%d", i),
		}
	}
}

var (
	auBcastChannelMap = channelMap{
		{5, 12, 177500000, 7000000},
		{21, 69, 480500000, 7000000},
	}

	euCableChannelMap = channelMap{
		{108, 862, 108000000, 1000000},
	}

	euBcastChannelMap = channelMap{
		{5, 12, 177500000, 7000000},
		{21, 69, 474000000, 8000000},
	}

	krCableChannelMap = channelMap{
		{2, 4, 57000000, 6000000},
		{5, 6, 79000000, 6000000},
		{7, 13, 177000000, 6000000},
		{14, 22, 123000000, 6000000},
		{23, 153, 219000000, 6000000},
	}

	usBcastChannelMap = channelMap{
		{2, 4, 57000000, 6000000},
		{5, 6, 79000000, 6000000},
		{7, 13, 177000000, 6000000},
		{14, 69, 473000000, 6000000},
	}

	usCableChannelMap = channelMap{
		{2, 4, 57000000, 6000000},
		{5, 6, 79000000, 6000000},
		{7, 13, 177000000, 6000000},
		{14, 22, 123000000, 6000000},
		{23, 94, 219000000, 6000000},
		{95, 99, 93000000, 6000000},
		{100, 158, 651000000, 6000000},
	}

	usHrcChannelMap = channelMap{
		{2, 4, 55752700, 6000300},
		{5, 6, 79753900, 6000300},
		{7, 13, 175758700, 6000300},
		{14, 22, 121756000, 6000300},
		{23, 94, 217760800, 6000300},
		{95, 99, 91754500, 6000300},
		{100, 158, 649782400, 6000300},
	}

	usIrcChannelMap = channelMap{
		{2, 4, 57012500, 6000000},
		{5, 6, 81012500, 6000000},
		{7, 13, 177012500, 6000000},
		{14, 22, 123012500, 6000000},
		{23, 41, 219012500, 6000000},
		{42, 42, 333025000, 6000000},
		{43, 94, 339012500, 6000000},
		{95, 97, 93012500, 6000000},
		{98, 99, 111025000, 6000000},
		{100, 158, 651012500, 6000000},
	}

	jpBcastChannelMap = channelMap{
		{13, 62, 473000000, 6000000},
	}

	channelMapTable = map[string]channelMap{
		"au-bcast": auBcastChannelMap,
		"au-cable": euCableChannelMap,
		"eu-bcast": euBcastChannelMap,
		"eu-cable": euCableChannelMap,
		"tw-bcast": usBcastChannelMap,
		"tw-cable": usCableChannelMap,
		"kr-bcast": usBcastChannelMap,
		"kr-cable": krCableChannelMap,
		"us-bcast": usBcastChannelMap,
		"us-cable": append(usCableChannelMap, append(usHrcChannelMap, usIrcChannelMap...)...),
		"us-hrc":   append(usCableChannelMap, append(usHrcChannelMap, usIrcChannelMap...)...),
		"us-irc":   append(usCableChannelMap, append(usHrcChannelMap, usIrcChannelMap...)...),
		"jp-bcast": jpBcastChannelMap,
	}
)

func Channels(channelMap string) chan Channel {
	ch := make(chan Channel)

	go func() {
		var wg sync.WaitGroup

		for _, cr := range channelMapTable[channelMap] {
			wg.Add(1)
			go func(cr channelRange) { defer wg.Done(); cr.channels(ch) }(cr)
		}

		wg.Wait()
		close(ch)
	}()

	return ch
}
