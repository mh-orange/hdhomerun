package hdhomerun

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
)

type ProgramType int

const (
	ProgramTypeNormal = iota
	ProgramTypeNoData
	ProgramTypeControl
	ProgramTypeEncrypted
)

func (pt ProgramType) String() string {
	switch pt {
	case ProgramTypeNormal:
		return "normal"
	case ProgramTypeNoData:
		return "no data"
	case ProgramTypeControl:
		return "control"
	case ProgramTypeEncrypted:
		return "encrypted"
	default:
		return "unknown"
	}
}

type Program struct {
	Name         string
	Type         ProgramType
	Number       int
	VirtualMajor int
	VirtualMinor int
}

func (p *Program) UnmarshalText(value TagValue) (err error) {
	str := string(value)
	tokens := strings.Split(str, ": ")
	if len(tokens) < 2 {
		return ErrParseError("Failed to parse " + str)
	}
	p.Number, err = parseInt(tokens[0])
	if err != nil {
		return err
	}
	tokens = strings.Split(tokens[1], " ")
	i := strings.Index(tokens[0], ".")
	if i != -1 {
		p.VirtualMajor, err = parseInt(tokens[0][0:i])
		if err == nil {
			p.VirtualMinor, err = parseInt(tokens[0][i+1:])
		}
	} else {
		p.VirtualMajor, err = parseInt(tokens[0])
	}

	if err != nil {
		return err
	}
	tokens = tokens[1:]

	if len(tokens) == 0 {
		return ErrParseError("No program name or tags")
	}

	if tokens[len(tokens)-1] == "(control)" {
		p.Type = ProgramTypeControl
		tokens = tokens[0 : len(tokens)-1]
	} else if tokens[len(tokens)-1] == "(encrypted)" {
		p.Type = ProgramTypeEncrypted
		tokens = tokens[0 : len(tokens)-1]
	} else if len(tokens) >= 2 && strings.Join(tokens[len(tokens)-2:], " ") == "(no data)" {
		p.Type = ProgramTypeNoData
		tokens = tokens[0 : len(tokens)-2]
	} else {
		p.Type = ProgramTypeNormal
	}

	p.Name = strings.Join(tokens, " ")

	return nil
}

func (p *Program) MarshalText() ([]byte, error) {
	str := ""
	if p.VirtualMinor != 0 {
		str = fmt.Sprintf("%d: %d.%d %s", p.Number, p.VirtualMajor, p.VirtualMinor, p.Name)
	} else {
		str = fmt.Sprintf("%d: %d %s", p.Number, p.VirtualMajor, p.Name)
	}

	if p.Type != ProgramTypeNormal {
		str += fmt.Sprintf(" (%s)", p.Type.String())
	}
	return []byte(str), nil
}

type ProgramList []Program

type Channel struct {
	Number     int
	Frequency  uint32
	Modulation string
	Name       string
	TSID       int
	ONID       int
	Programs   []Program
}

func (c *Channel) UnmarshalText(text []byte) (err error) {
	reader := bufio.NewReader(bytes.NewBuffer(text))
	var line []byte
	var isPrefix bool
	for err == nil {
		line, isPrefix, err = reader.ReadLine()
		if err != nil {
			break
		}

		if isPrefix {
			return ErrParseError("The program line was too long")
		}

		n := 0
		if n, _ = fmt.Sscanf(string(line), "tsid=0x%x", &c.TSID); n == 1 {
			continue
		}

		if n, _ = fmt.Sscanf(string(line), "onid=0x%x", &c.ONID); n == 1 {
			continue
		}

		p := &Program{}
		err = p.UnmarshalText(line)
		if err == nil {
			c.Programs = append(c.Programs, *p)
		}
	}

	if err == io.EOF {
		err = nil
	}
	return
}

func (c *Channel) TSIDDetected() bool {
	return c.TSID >= 0
}

func (c *Channel) ONIDDetected() bool {
	return c.ONID >= 0
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
			TSID:       -1,
			ONID:       -1,
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
