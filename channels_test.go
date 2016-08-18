package hdhomerun

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestChannels(t *testing.T) {
	tests := []struct {
		cr       channelRange
		expected []Channel
	}{
		{
			cr: channelRange{
				Start:     2,
				End:       4,
				Frequency: 57000000,
				Spacing:   6000000,
			},
			expected: []Channel{
				Channel{2, 57000000, "auto", "2", -1, -1, nil},
				Channel{3, 63000000, "auto", "3", -1, -1, nil},
				Channel{4, 69000000, "auto", "4", -1, -1, nil},
			},
		},
		{
			cr: channelRange{
				Start:     7,
				End:       13,
				Frequency: 177000000,
				Spacing:   6000000,
			},
			expected: []Channel{
				Channel{7, 177000000, "auto", "7", -1, -1, nil},
				Channel{8, 183000000, "auto", "8", -1, -1, nil},
				Channel{9, 189000000, "auto", "9", -1, -1, nil},
				Channel{10, 195000000, "auto", "10", -1, -1, nil},
				Channel{11, 201000000, "auto", "11", -1, -1, nil},
				Channel{12, 207000000, "auto", "12", -1, -1, nil},
				Channel{13, 213000000, "auto", "13", -1, -1, nil},
			},
		},
	}

	for _, test := range tests {
		received := test.cr.Channels()

		if !reflect.DeepEqual(test.expected, received) {
			t.Errorf("Expected \n%v\ngot\n%v\n", test.expected, received)
		}
	}
}

func TestChannelMap(t *testing.T) {
	channels := channelMap{
		{2, 4, 57000000, 6000000},
		{5, 6, 79000000, 6000000},
		{7, 13, 177000000, 6000000},
		{14, 22, 123000000, 6000000},
		{23, 94, 219000000, 6000000},
		{95, 99, 93000000, 6000000},
		{100, 158, 651000000, 6000000},
		{2, 4, 55752700, 6000300},
		{5, 6, 79753900, 6000300},
		{7, 13, 175758700, 6000300},
		{14, 22, 121756000, 6000300},
		{23, 94, 217760800, 6000300},
		{95, 99, 91754500, 6000300},
		{100, 158, 649782400, 6000300},
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

	if !reflect.DeepEqual(channels, channelMapTable["us-cable"]) {
		t.Errorf("Expected\n%v\nGot\n%v\n", channels, channelMapTable["us-cable"])
	}
}

func TestChannelMapChannels(t *testing.T) {

	channels := []Channel{
		Channel{97, 103750000, "auto", "97", -1, -1, nil},
		Channel{97, 105000000, "auto", "97", -1, -1, nil},
		Channel{98, 111000000, "auto", "98", -1, -1, nil},
		Channel{99, 117000000, "auto", "99", -1, -1, nil},
		Channel{14, 123000000, "auto", "14", -1, -1, nil},
		Channel{15, 129000000, "auto", "15", -1, -1, nil},
		Channel{16, 135000000, "auto", "16", -1, -1, nil},
		Channel{17, 141000000, "auto", "17", -1, -1, nil},
		Channel{18, 147000000, "auto", "18", -1, -1, nil},
		Channel{19, 153000000, "auto", "19", -1, -1, nil},
		Channel{20, 159000000, "auto", "20", -1, -1, nil},
		Channel{21, 165000000, "auto", "21", -1, -1, nil},
		Channel{22, 171000000, "auto", "22", -1, -1, nil},
		Channel{7, 177000000, "auto", "7", -1, -1, nil},
		Channel{8, 183000000, "auto", "8", -1, -1, nil},
		Channel{9, 189000000, "auto", "9", -1, -1, nil},
		Channel{10, 195000000, "auto", "10", -1, -1, nil},
		Channel{11, 201000000, "auto", "11", -1, -1, nil},
		Channel{12, 207000000, "auto", "12", -1, -1, nil},
		Channel{13, 213000000, "auto", "13", -1, -1, nil},
		Channel{23, 219000000, "auto", "23", -1, -1, nil},
		Channel{24, 225000000, "auto", "24", -1, -1, nil},
		Channel{25, 231000000, "auto", "25", -1, -1, nil},
		Channel{26, 237000000, "auto", "26", -1, -1, nil},
		Channel{27, 243000000, "auto", "27", -1, -1, nil},
		Channel{28, 249000000, "auto", "28", -1, -1, nil},
		Channel{29, 255000000, "auto", "29", -1, -1, nil},
		Channel{30, 261000000, "auto", "30", -1, -1, nil},
		Channel{31, 267000000, "auto", "31", -1, -1, nil},
		Channel{32, 273000000, "auto", "32", -1, -1, nil},
		Channel{33, 279000000, "auto", "33", -1, -1, nil},
		Channel{34, 285000000, "auto", "34", -1, -1, nil},
		Channel{35, 291000000, "auto", "35", -1, -1, nil},
		Channel{36, 297000000, "auto", "36", -1, -1, nil},
		Channel{37, 303000000, "auto", "37", -1, -1, nil},
		Channel{38, 309000000, "auto", "38", -1, -1, nil},
		Channel{39, 315000000, "auto", "39", -1, -1, nil},
		Channel{40, 321000000, "auto", "40", -1, -1, nil},
		Channel{41, 327000000, "auto", "41", -1, -1, nil},
		Channel{42, 333000000, "auto", "42", -1, -1, nil},
		Channel{43, 339000000, "auto", "43", -1, -1, nil},
		Channel{44, 345000000, "auto", "44", -1, -1, nil},
		Channel{45, 351000000, "auto", "45", -1, -1, nil},
		Channel{46, 357000000, "auto", "46", -1, -1, nil},
		Channel{47, 363000000, "auto", "47", -1, -1, nil},
		Channel{48, 369000000, "auto", "48", -1, -1, nil},
		Channel{49, 375000000, "auto", "49", -1, -1, nil},
		Channel{50, 381000000, "auto", "50", -1, -1, nil},
		Channel{51, 387000000, "auto", "51", -1, -1, nil},
		Channel{52, 393000000, "auto", "52", -1, -1, nil},
		Channel{53, 399000000, "auto", "53", -1, -1, nil},
		Channel{54, 405000000, "auto", "54", -1, -1, nil},
		Channel{55, 411000000, "auto", "55", -1, -1, nil},
		Channel{56, 417000000, "auto", "56", -1, -1, nil},
		Channel{57, 423000000, "auto", "57", -1, -1, nil},
		Channel{58, 429000000, "auto", "58", -1, -1, nil},
		Channel{59, 435000000, "auto", "59", -1, -1, nil},
		Channel{60, 441000000, "auto", "60", -1, -1, nil},
		Channel{61, 447000000, "auto", "61", -1, -1, nil},
		Channel{62, 453000000, "auto", "62", -1, -1, nil},
		Channel{63, 459000000, "auto", "63", -1, -1, nil},
		Channel{64, 465000000, "auto", "64", -1, -1, nil},
		Channel{65, 471000000, "auto", "65", -1, -1, nil},
		Channel{66, 477000000, "auto", "66", -1, -1, nil},
		Channel{67, 483000000, "auto", "67", -1, -1, nil},
		Channel{68, 489000000, "auto", "68", -1, -1, nil},
		Channel{69, 495000000, "auto", "69", -1, -1, nil},
		Channel{70, 501000000, "auto", "70", -1, -1, nil},
		Channel{71, 507000000, "auto", "71", -1, -1, nil},
		Channel{72, 513000000, "auto", "72", -1, -1, nil},
		Channel{73, 519000000, "auto", "73", -1, -1, nil},
		Channel{74, 525000000, "auto", "74", -1, -1, nil},
		Channel{75, 531000000, "auto", "75", -1, -1, nil},
		Channel{76, 537000000, "auto", "76", -1, -1, nil},
		Channel{77, 543000000, "auto", "77", -1, -1, nil},
		Channel{78, 549000000, "auto", "78", -1, -1, nil},
		Channel{79, 555000000, "auto", "79", -1, -1, nil},
		Channel{80, 561000000, "auto", "80", -1, -1, nil},
		Channel{81, 567000000, "auto", "81", -1, -1, nil},
		Channel{2, 57000000, "auto", "2", -1, -1, nil},
		Channel{82, 573000000, "auto", "82", -1, -1, nil},
		Channel{83, 579000000, "auto", "83", -1, -1, nil},
		Channel{84, 585000000, "auto", "84", -1, -1, nil},
		Channel{85, 591000000, "auto", "85", -1, -1, nil},
		Channel{86, 597000000, "auto", "86", -1, -1, nil},
		Channel{87, 603000000, "auto", "87", -1, -1, nil},
		Channel{88, 609000000, "auto", "88", -1, -1, nil},
		Channel{89, 615000000, "auto", "89", -1, -1, nil},
		Channel{90, 621000000, "auto", "90", -1, -1, nil},
		Channel{91, 627000000, "auto", "91", -1, -1, nil},
		Channel{3, 63000000, "auto", "3", -1, -1, nil},
		Channel{92, 633000000, "auto", "92", -1, -1, nil},
		Channel{93, 639000000, "auto", "93", -1, -1, nil},
		Channel{94, 645000000, "auto", "94", -1, -1, nil},
		Channel{100, 651000000, "auto", "100", -1, -1, nil},
		Channel{101, 657000000, "auto", "101", -1, -1, nil},
		Channel{102, 663000000, "auto", "102", -1, -1, nil},
		Channel{103, 669000000, "auto", "103", -1, -1, nil},
		Channel{104, 675000000, "auto", "104", -1, -1, nil},
		Channel{105, 681000000, "auto", "105", -1, -1, nil},
		Channel{106, 687000000, "auto", "106", -1, -1, nil},
		Channel{4, 69000000, "auto", "4", -1, -1, nil},
		Channel{107, 693000000, "auto", "107", -1, -1, nil},
		Channel{108, 699000000, "auto", "108", -1, -1, nil},
		Channel{109, 705000000, "auto", "109", -1, -1, nil},
		Channel{110, 711000000, "auto", "110", -1, -1, nil},
		Channel{111, 717000000, "auto", "111", -1, -1, nil},
		Channel{112, 723000000, "auto", "112", -1, -1, nil},
		Channel{113, 729000000, "auto", "113", -1, -1, nil},
		Channel{114, 735000000, "auto", "114", -1, -1, nil},
		Channel{115, 741000000, "auto", "115", -1, -1, nil},
		Channel{116, 747000000, "auto", "116", -1, -1, nil},
		Channel{117, 753000000, "auto", "117", -1, -1, nil},
		Channel{118, 759000000, "auto", "118", -1, -1, nil},
		Channel{119, 765000000, "auto", "119", -1, -1, nil},
		Channel{120, 771000000, "auto", "120", -1, -1, nil},
		Channel{121, 777000000, "auto", "121", -1, -1, nil},
		Channel{122, 783000000, "auto", "122", -1, -1, nil},
		Channel{123, 789000000, "auto", "123", -1, -1, nil},
		Channel{5, 79000000, "auto", "5", -1, -1, nil},
		Channel{124, 795000000, "auto", "124", -1, -1, nil},
		Channel{125, 801000000, "auto", "125", -1, -1, nil},
		Channel{126, 807000000, "auto", "126", -1, -1, nil},
		Channel{127, 813000000, "auto", "127", -1, -1, nil},
		Channel{128, 819000000, "auto", "128", -1, -1, nil},
		Channel{129, 823750000, "auto", "129", -1, -1, nil},
		Channel{129, 825000000, "auto", "129", -1, -1, nil},
		Channel{130, 831000000, "auto", "130", -1, -1, nil},
		Channel{131, 837000000, "auto", "131", -1, -1, nil},
		Channel{132, 843000000, "auto", "132", -1, -1, nil},
		Channel{133, 849000000, "auto", "133", -1, -1, nil},
		Channel{6, 85000000, "auto", "6", -1, -1, nil},
		Channel{134, 855000000, "auto", "134", -1, -1, nil},
		Channel{6, 85750000, "auto", "6", -1, -1, nil},
		Channel{135, 861000000, "auto", "135", -1, -1, nil},
		Channel{136, 865750000, "auto", "136", -1, -1, nil},
		Channel{136, 867000000, "auto", "136", -1, -1, nil},
		Channel{6, 87000000, "auto", "6", -1, -1, nil},
		Channel{137, 871750000, "auto", "137", -1, -1, nil},
		Channel{137, 873000000, "auto", "137", -1, -1, nil},
		Channel{138, 877750000, "auto", "138", -1, -1, nil},
		Channel{138, 879000000, "auto", "138", -1, -1, nil},
		Channel{139, 883750000, "auto", "139", -1, -1, nil},
		Channel{139, 885000000, "auto", "139", -1, -1, nil},
		Channel{140, 889750000, "auto", "140", -1, -1, nil},
		Channel{140, 891000000, "auto", "140", -1, -1, nil},
		Channel{141, 895750000, "auto", "141", -1, -1, nil},
		Channel{141, 897000000, "auto", "141", -1, -1, nil},
		Channel{142, 901750000, "auto", "142", -1, -1, nil},
		Channel{142, 903000000, "auto", "142", -1, -1, nil},
		Channel{143, 907750000, "auto", "143", -1, -1, nil},
		Channel{143, 909000000, "auto", "143", -1, -1, nil},
		Channel{144, 913750000, "auto", "144", -1, -1, nil},
		Channel{144, 915000000, "auto", "144", -1, -1, nil},
		Channel{95, 91750000, "auto", "95", -1, -1, nil},
		Channel{145, 919750000, "auto", "145", -1, -1, nil},
		Channel{145, 921000000, "auto", "145", -1, -1, nil},
		Channel{146, 925750000, "auto", "146", -1, -1, nil},
		Channel{146, 927000000, "auto", "146", -1, -1, nil},
		Channel{95, 93000000, "auto", "95", -1, -1, nil},
		Channel{147, 931750000, "auto", "147", -1, -1, nil},
		Channel{147, 933000000, "auto", "147", -1, -1, nil},
		Channel{148, 937750000, "auto", "148", -1, -1, nil},
		Channel{148, 939000000, "auto", "148", -1, -1, nil},
		Channel{149, 943750000, "auto", "149", -1, -1, nil},
		Channel{149, 945000000, "auto", "149", -1, -1, nil},
		Channel{150, 949750000, "auto", "150", -1, -1, nil},
		Channel{150, 951000000, "auto", "150", -1, -1, nil},
		Channel{151, 955750000, "auto", "151", -1, -1, nil},
		Channel{151, 957000000, "auto", "151", -1, -1, nil},
		Channel{152, 961750000, "auto", "152", -1, -1, nil},
		Channel{152, 963000000, "auto", "152", -1, -1, nil},
		Channel{153, 967750000, "auto", "153", -1, -1, nil},
		Channel{153, 969000000, "auto", "153", -1, -1, nil},
		Channel{154, 973750000, "auto", "154", -1, -1, nil},
		Channel{154, 975000000, "auto", "154", -1, -1, nil},
		Channel{96, 97750000, "auto", "96", -1, -1, nil},
		Channel{155, 979750000, "auto", "155", -1, -1, nil},
		Channel{155, 981000000, "auto", "155", -1, -1, nil},
		Channel{156, 985750000, "auto", "156", -1, -1, nil},
		Channel{156, 987000000, "auto", "156", -1, -1, nil},
		Channel{96, 99000000, "auto", "96", -1, -1, nil},
		Channel{157, 991750000, "auto", "157", -1, -1, nil},
		Channel{157, 993000000, "auto", "157", -1, -1, nil},
		Channel{158, 997750000, "auto", "158", -1, -1, nil},
		Channel{158, 999000000, "auto", "158", -1, -1, nil},
	}

	expected := make(map[int]Channel)
	for _, channel := range channels {
		expected[channel.Frequency] = channel
	}

	received := make(map[int]Channel)
	for channel := range Channels("us-cable") {
		received[channel.Frequency] = channel
	}

	for frequency, expChannel := range expected {
		if rxChannel, found := received[frequency]; found {
			if !reflect.DeepEqual(expChannel, rxChannel) {
				t.Errorf("Expected\n%v\nto be\n%v\n", rxChannel, expChannel)
			}
		} else {
			t.Errorf("Expected to receive channel %v but didn't get it", expChannel)
		}
	}

	for frequency, _ := range received {
		if expChannel, found := received[frequency]; !found {
			t.Errorf("Did not expecte to get channel %v", expChannel)
		}
	}

}

func TestMarshallingProgram(t *testing.T) {
	tests := []struct {
		str         string
		expected    *Program
		expectedErr reflect.Type
	}{
		{
			str: "735: 695 Comedy.TV HD (encrypted)",
			expected: &Program{
				Name:         "Comedy.TV HD",
				Type:         ProgramTypeEncrypted,
				Number:       735,
				VirtualMajor: 695,
				VirtualMinor: 0,
			},
		}, {
			str: "735: 695.5 Comedy.TV HD (encrypted)",
			expected: &Program{
				Name:         "Comedy.TV HD",
				Type:         ProgramTypeEncrypted,
				Number:       735,
				VirtualMajor: 695,
				VirtualMinor: 5,
			},
		}, {
			str: "735: 695 Comedy.TV HD (It's Great!) (encrypted)",
			expected: &Program{
				Name:         "Comedy.TV HD (It's Great!)",
				Type:         ProgramTypeEncrypted,
				Number:       735,
				VirtualMajor: 695,
				VirtualMinor: 0,
			},
		}, {
			str: "735: 695 Comedy.TV HD (It's Great!) (control)",
			expected: &Program{
				Name:         "Comedy.TV HD (It's Great!)",
				Type:         ProgramTypeControl,
				Number:       735,
				VirtualMajor: 695,
				VirtualMinor: 0,
			},
		}, {
			str: "735: 695 Comedy.TV HD (It's Great!) (no data)",
			expected: &Program{
				Name:         "Comedy.TV HD (It's Great!)",
				Type:         ProgramTypeNoData,
				Number:       735,
				VirtualMajor: 695,
				VirtualMinor: 0,
			},
		}, {
			str: "735: 695 Comedy.TV HD (It's Great!)",
			expected: &Program{
				Name:         "Comedy.TV HD (It's Great!)",
				Type:         ProgramTypeNormal,
				Number:       735,
				VirtualMajor: 695,
				VirtualMinor: 0,
			},
		}, {
			str:         "",
			expected:    nil,
			expectedErr: reflect.TypeOf(ErrParseError("")),
		}, {
			str:         "a: 695",
			expected:    nil,
			expectedErr: reflect.TypeOf(&strconv.NumError{}),
		}, {
			str:         "1: b",
			expected:    nil,
			expectedErr: reflect.TypeOf(&strconv.NumError{}),
		}, {
			str:         "1: 695",
			expected:    nil,
			expectedErr: reflect.TypeOf(ErrParseError("")),
		},
	}

	for _, test := range tests {
		p := &Program{}
		err := p.UnmarshalText([]byte(test.str))
		if err != nil {
			if reflect.TypeOf(err) != test.expectedErr {
				t.Errorf("Expected %v but got %v", test.expectedErr, err)
			}
			continue
		}

		if !reflect.DeepEqual(p, test.expected) {
			t.Errorf("Expected\n%v\nReceived:\n%v\n", test.expected, p)
		}

		b, _ := p.MarshalText()
		if !reflect.DeepEqual([]byte(test.str), b) {
			t.Errorf("Marshaling failed.  Expected:\n%s\nReceived\n%s\n", test.str, string(b))
		}
	}
}

func TestUnmarshalChannel(t *testing.T) {
	tests := []struct {
		str          string
		expected     *Channel
		expectedErr  reflect.Type
		tsidDetected bool
		onidDetected bool
	}{
		{
			str: "735: 695 Comedy.TV HD (encrypted)\n736: 599 Cars.TV HD (encrypted)\ntsid=0x0001\n",
			expected: &Channel{
				TSID: 1,
				ONID: -1,
				Programs: []Program{
					Program{
						Number:       735,
						VirtualMajor: 695,
						Name:         "Comedy.TV HD",
						Type:         ProgramTypeEncrypted,
					},
					Program{
						Number:       736,
						VirtualMajor: 599,
						Name:         "Cars.TV HD",
						Type:         ProgramTypeEncrypted,
					},
				},
			},
			tsidDetected: true,
		}, {
			str: "onid=0x0001\n",
			expected: &Channel{
				TSID:     -1,
				ONID:     1,
				Programs: nil,
			},
			onidDetected: true,
		}, {
			str: strings.Repeat("S", 4097),
			expected: &Channel{
				TSID:     -1,
				ONID:     -1,
				Programs: nil,
			},
			expectedErr: reflect.TypeOf(ErrParseError("")),
		},
	}

	for _, test := range tests {
		c := &Channel{
			ONID: -1,
			TSID: -1,
		}

		err := c.UnmarshalText([]byte(test.str))
		if err != nil {
			if reflect.TypeOf(err) != test.expectedErr {
				t.Errorf("Expected %v but got %v", test.expectedErr, err)
			}
			continue
		}

		if !reflect.DeepEqual(c, test.expected) {
			t.Errorf("Expected\n%v\nReceived:\n%v\n", test.expected, c)
		}

		if c.TSIDDetected() != test.tsidDetected {
			t.Errorf("Expected TSIDDetected to be %v but got %v", test.tsidDetected, c.TSIDDetected())
		}

		if c.ONIDDetected() != test.onidDetected {
			t.Errorf("Expected ONIDDetected to be %v but got %v", test.onidDetected, c.ONIDDetected())
		}
	}
}

func TestChannelTypeString(t *testing.T) {
	tests := []struct {
		t        ProgramType
		expected string
	}{
		{
			t:        ProgramTypeNormal,
			expected: "normal",
		}, {
			t:        ProgramTypeNoData,
			expected: "no data",
		}, {
			t:        ProgramTypeControl,
			expected: "control",
		}, {
			t:        ProgramTypeEncrypted,
			expected: "encrypted",
		}, {
			t:        ProgramType(-1),
			expected: "unknown",
		},
	}

	for _, test := range tests {
		if test.t.String() != test.expected {
			t.Errorf("Expected %s got %s", test.expected, test.t.String())
		}
	}
}
