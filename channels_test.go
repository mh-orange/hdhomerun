package hdhomerun

import (
	"reflect"
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
				Channel{2, 57000000, "auto", "2", nil},
				Channel{3, 63000000, "auto", "3", nil},
				Channel{4, 69000000, "auto", "4", nil},
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
				Channel{7, 177000000, "auto", "7", nil},
				Channel{8, 183000000, "auto", "8", nil},
				Channel{9, 189000000, "auto", "9", nil},
				Channel{10, 195000000, "auto", "10", nil},
				Channel{11, 201000000, "auto", "11", nil},
				Channel{12, 207000000, "auto", "12", nil},
				Channel{13, 213000000, "auto", "13", nil},
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
		Channel{97, 103750000, "auto", "97", nil},
		Channel{97, 105000000, "auto", "97", nil},
		Channel{98, 111000000, "auto", "98", nil},
		Channel{99, 117000000, "auto", "99", nil},
		Channel{14, 123000000, "auto", "14", nil},
		Channel{15, 129000000, "auto", "15", nil},
		Channel{16, 135000000, "auto", "16", nil},
		Channel{17, 141000000, "auto", "17", nil},
		Channel{18, 147000000, "auto", "18", nil},
		Channel{19, 153000000, "auto", "19", nil},
		Channel{20, 159000000, "auto", "20", nil},
		Channel{21, 165000000, "auto", "21", nil},
		Channel{22, 171000000, "auto", "22", nil},
		Channel{7, 177000000, "auto", "7", nil},
		Channel{8, 183000000, "auto", "8", nil},
		Channel{9, 189000000, "auto", "9", nil},
		Channel{10, 195000000, "auto", "10", nil},
		Channel{11, 201000000, "auto", "11", nil},
		Channel{12, 207000000, "auto", "12", nil},
		Channel{13, 213000000, "auto", "13", nil},
		Channel{23, 219000000, "auto", "23", nil},
		Channel{24, 225000000, "auto", "24", nil},
		Channel{25, 231000000, "auto", "25", nil},
		Channel{26, 237000000, "auto", "26", nil},
		Channel{27, 243000000, "auto", "27", nil},
		Channel{28, 249000000, "auto", "28", nil},
		Channel{29, 255000000, "auto", "29", nil},
		Channel{30, 261000000, "auto", "30", nil},
		Channel{31, 267000000, "auto", "31", nil},
		Channel{32, 273000000, "auto", "32", nil},
		Channel{33, 279000000, "auto", "33", nil},
		Channel{34, 285000000, "auto", "34", nil},
		Channel{35, 291000000, "auto", "35", nil},
		Channel{36, 297000000, "auto", "36", nil},
		Channel{37, 303000000, "auto", "37", nil},
		Channel{38, 309000000, "auto", "38", nil},
		Channel{39, 315000000, "auto", "39", nil},
		Channel{40, 321000000, "auto", "40", nil},
		Channel{41, 327000000, "auto", "41", nil},
		Channel{42, 333000000, "auto", "42", nil},
		Channel{43, 339000000, "auto", "43", nil},
		Channel{44, 345000000, "auto", "44", nil},
		Channel{45, 351000000, "auto", "45", nil},
		Channel{46, 357000000, "auto", "46", nil},
		Channel{47, 363000000, "auto", "47", nil},
		Channel{48, 369000000, "auto", "48", nil},
		Channel{49, 375000000, "auto", "49", nil},
		Channel{50, 381000000, "auto", "50", nil},
		Channel{51, 387000000, "auto", "51", nil},
		Channel{52, 393000000, "auto", "52", nil},
		Channel{53, 399000000, "auto", "53", nil},
		Channel{54, 405000000, "auto", "54", nil},
		Channel{55, 411000000, "auto", "55", nil},
		Channel{56, 417000000, "auto", "56", nil},
		Channel{57, 423000000, "auto", "57", nil},
		Channel{58, 429000000, "auto", "58", nil},
		Channel{59, 435000000, "auto", "59", nil},
		Channel{60, 441000000, "auto", "60", nil},
		Channel{61, 447000000, "auto", "61", nil},
		Channel{62, 453000000, "auto", "62", nil},
		Channel{63, 459000000, "auto", "63", nil},
		Channel{64, 465000000, "auto", "64", nil},
		Channel{65, 471000000, "auto", "65", nil},
		Channel{66, 477000000, "auto", "66", nil},
		Channel{67, 483000000, "auto", "67", nil},
		Channel{68, 489000000, "auto", "68", nil},
		Channel{69, 495000000, "auto", "69", nil},
		Channel{70, 501000000, "auto", "70", nil},
		Channel{71, 507000000, "auto", "71", nil},
		Channel{72, 513000000, "auto", "72", nil},
		Channel{73, 519000000, "auto", "73", nil},
		Channel{74, 525000000, "auto", "74", nil},
		Channel{75, 531000000, "auto", "75", nil},
		Channel{76, 537000000, "auto", "76", nil},
		Channel{77, 543000000, "auto", "77", nil},
		Channel{78, 549000000, "auto", "78", nil},
		Channel{79, 555000000, "auto", "79", nil},
		Channel{80, 561000000, "auto", "80", nil},
		Channel{81, 567000000, "auto", "81", nil},
		Channel{2, 57000000, "auto", "2", nil},
		Channel{82, 573000000, "auto", "82", nil},
		Channel{83, 579000000, "auto", "83", nil},
		Channel{84, 585000000, "auto", "84", nil},
		Channel{85, 591000000, "auto", "85", nil},
		Channel{86, 597000000, "auto", "86", nil},
		Channel{87, 603000000, "auto", "87", nil},
		Channel{88, 609000000, "auto", "88", nil},
		Channel{89, 615000000, "auto", "89", nil},
		Channel{90, 621000000, "auto", "90", nil},
		Channel{91, 627000000, "auto", "91", nil},
		Channel{3, 63000000, "auto", "3", nil},
		Channel{92, 633000000, "auto", "92", nil},
		Channel{93, 639000000, "auto", "93", nil},
		Channel{94, 645000000, "auto", "94", nil},
		Channel{100, 651000000, "auto", "100", nil},
		Channel{101, 657000000, "auto", "101", nil},
		Channel{102, 663000000, "auto", "102", nil},
		Channel{103, 669000000, "auto", "103", nil},
		Channel{104, 675000000, "auto", "104", nil},
		Channel{105, 681000000, "auto", "105", nil},
		Channel{106, 687000000, "auto", "106", nil},
		Channel{4, 69000000, "auto", "4", nil},
		Channel{107, 693000000, "auto", "107", nil},
		Channel{108, 699000000, "auto", "108", nil},
		Channel{109, 705000000, "auto", "109", nil},
		Channel{110, 711000000, "auto", "110", nil},
		Channel{111, 717000000, "auto", "111", nil},
		Channel{112, 723000000, "auto", "112", nil},
		Channel{113, 729000000, "auto", "113", nil},
		Channel{114, 735000000, "auto", "114", nil},
		Channel{115, 741000000, "auto", "115", nil},
		Channel{116, 747000000, "auto", "116", nil},
		Channel{117, 753000000, "auto", "117", nil},
		Channel{118, 759000000, "auto", "118", nil},
		Channel{119, 765000000, "auto", "119", nil},
		Channel{120, 771000000, "auto", "120", nil},
		Channel{121, 777000000, "auto", "121", nil},
		Channel{122, 783000000, "auto", "122", nil},
		Channel{123, 789000000, "auto", "123", nil},
		Channel{5, 79000000, "auto", "5", nil},
		Channel{124, 795000000, "auto", "124", nil},
		Channel{125, 801000000, "auto", "125", nil},
		Channel{126, 807000000, "auto", "126", nil},
		Channel{127, 813000000, "auto", "127", nil},
		Channel{128, 819000000, "auto", "128", nil},
		Channel{129, 823750000, "auto", "129", nil},
		Channel{129, 825000000, "auto", "129", nil},
		Channel{130, 831000000, "auto", "130", nil},
		Channel{131, 837000000, "auto", "131", nil},
		Channel{132, 843000000, "auto", "132", nil},
		Channel{133, 849000000, "auto", "133", nil},
		Channel{6, 85000000, "auto", "6", nil},
		Channel{134, 855000000, "auto", "134", nil},
		Channel{6, 85750000, "auto", "6", nil},
		Channel{135, 861000000, "auto", "135", nil},
		Channel{136, 865750000, "auto", "136", nil},
		Channel{136, 867000000, "auto", "136", nil},
		Channel{6, 87000000, "auto", "6", nil},
		Channel{137, 871750000, "auto", "137", nil},
		Channel{137, 873000000, "auto", "137", nil},
		Channel{138, 877750000, "auto", "138", nil},
		Channel{138, 879000000, "auto", "138", nil},
		Channel{139, 883750000, "auto", "139", nil},
		Channel{139, 885000000, "auto", "139", nil},
		Channel{140, 889750000, "auto", "140", nil},
		Channel{140, 891000000, "auto", "140", nil},
		Channel{141, 895750000, "auto", "141", nil},
		Channel{141, 897000000, "auto", "141", nil},
		Channel{142, 901750000, "auto", "142", nil},
		Channel{142, 903000000, "auto", "142", nil},
		Channel{143, 907750000, "auto", "143", nil},
		Channel{143, 909000000, "auto", "143", nil},
		Channel{144, 913750000, "auto", "144", nil},
		Channel{144, 915000000, "auto", "144", nil},
		Channel{95, 91750000, "auto", "95", nil},
		Channel{145, 919750000, "auto", "145", nil},
		Channel{145, 921000000, "auto", "145", nil},
		Channel{146, 925750000, "auto", "146", nil},
		Channel{146, 927000000, "auto", "146", nil},
		Channel{95, 93000000, "auto", "95", nil},
		Channel{147, 931750000, "auto", "147", nil},
		Channel{147, 933000000, "auto", "147", nil},
		Channel{148, 937750000, "auto", "148", nil},
		Channel{148, 939000000, "auto", "148", nil},
		Channel{149, 943750000, "auto", "149", nil},
		Channel{149, 945000000, "auto", "149", nil},
		Channel{150, 949750000, "auto", "150", nil},
		Channel{150, 951000000, "auto", "150", nil},
		Channel{151, 955750000, "auto", "151", nil},
		Channel{151, 957000000, "auto", "151", nil},
		Channel{152, 961750000, "auto", "152", nil},
		Channel{152, 963000000, "auto", "152", nil},
		Channel{153, 967750000, "auto", "153", nil},
		Channel{153, 969000000, "auto", "153", nil},
		Channel{154, 973750000, "auto", "154", nil},
		Channel{154, 975000000, "auto", "154", nil},
		Channel{96, 97750000, "auto", "96", nil},
		Channel{155, 979750000, "auto", "155", nil},
		Channel{155, 981000000, "auto", "155", nil},
		Channel{156, 985750000, "auto", "156", nil},
		Channel{156, 987000000, "auto", "156", nil},
		Channel{96, 99000000, "auto", "96", nil},
		Channel{157, 991750000, "auto", "157", nil},
		Channel{157, 993000000, "auto", "157", nil},
		Channel{158, 997750000, "auto", "158", nil},
		Channel{158, 999000000, "auto", "158", nil},
	}

	expected := make(map[uint32]Channel)
	for _, channel := range channels {
		expected[channel.Frequency] = channel
	}

	received := make(map[uint32]Channel)
	for channel := range Channels("us-cable") {
		received[channel.Frequency] = channel
	}
	/*for _, cr := range ChannelMapTable["us-cable"] {
		for _, channel := range cr.Channels() {
			received[channel.Frequency] = channel
		}
	}*/

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
