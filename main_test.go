package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestBuildLookupTable(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      []tag
		expectedError error
	}{
		{"1",
			"25,tcp,sv_P1",
			[]tag{
				{dstPort: 25, protocol: 6, name: "sv_p1"},
			},
			nil,
		},
		{"2",
			"# comment",
			[]tag{},
			nil,
		},
		{"3",
			"1,2,3,4,5",
			nil,
			errors.New("record on line 1: wrong number of fields"),
		},
		{"4",
			"65536,tcp,sv_P1",
			nil,
			errors.New("invalid port 65536"),
		},
		{"5",
			"22,bla,sv_P1",
			nil,
			errors.New("invalid protocol number: bla"),
		},
		{"6",
			"string,tcp,sv_P1",
			nil,
			errors.New("strconv.Atoi: parsing \"string\": invalid syntax"),
		},
		{"7",
			`# comment
25,tcp,sv_P1
68,udp,sv_P2
`,
			[]tag{
				{dstPort: 25, protocol: 6, name: "sv_p1"},
				{dstPort: 68, protocol: 17, name: "sv_p2"},
			},
			nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tags, e := buildLookupTable(strings.NewReader(test.input))
			require.Equal(t, test.expected, tags)
			if test.expectedError != nil {
				require.EqualError(t, e, test.expectedError.Error())
			}
		})
	}
}

func TestStat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		tags     []tag
		m        map[string]int
		untagged int
	}{
		{"1",
			`
2 123456789012 eni-0a1b2c3d 10.0.1.201 198.51.100.2 443 49153 6 25 20000 1620140761 1620140821 ACCEPT OK
2 123456789012 eni-4d3c2b1a 192.168.1.100 203.0.113.101 23 49154 6 15 12000 1620140761 1620140821 REJECT OK
2 123456789012 eni-5e6f7g8h 192.168.1.101 198.51.100.3 25 49155 6 10 8000 1620140761 1620140821 ACCEPT OK
2 123456789012 eni-9h8g7f6e 172.16.0.100 203.0.113.102 110 49156 6 12 9000 1620140761 1620140821 ACCEPT OK
2 123456789012 eni-7i8j9k0l 172.16.0.101 192.0.2.203 993 49157 6 8 5000 1620140761 1620140821 ACCEPT OK
2 123456789012 eni-6m7n8o9p 10.0.2.200 198.51.100.4 143 49158 6 18 14000 1620140761 1620140821 ACCEPT OK
2 123456789012 eni-1a2b3c4d 192.168.0.1 203.0.113.12 1024 80 6 10 5000 1620140661 1620140721 ACCEPT OK
2 123456789012 eni-1a2b3c4d 203.0.113.12 192.168.0.1 80 1024 6 12 6000 1620140661 1620140721 ACCEPT OK
2 123456789012 eni-1a2b3c4d 10.0.1.102 172.217.7.228 1030 443 6 8 4000 1620140661 1620140721 ACCEPT OK
2 123456789012 eni-5f6g7h8i 10.0.2.103 52.26.198.183 56000 23 6 15 7500 1620140661 1620140721 REJECT OK
2 123456789012 eni-9k10l11m 192.168.1.5 51.15.99.115 49321 25 6 20 10000 1620140661 1620140721 ACCEPT OK
2 123456789012 eni-1a2b3c4d 192.168.1.6 87.250.250.242 49152 110 6 5 2500 1620140661 1620140721 ACCEPT OK
2 123456789012 eni-2d2e2f3g 192.168.2.7 77.88.55.80 49153 993 6 7 3500 1620140661 1620140721 ACCEPT OK
2 123456789012 eni-4h5i6j7k 172.16.0.2 192.0.2.146 49154 143 6 9 4500 1620140661 1620140721 ACCEPT OK`,
			[]tag{
				{dstPort: 25, protocol: 6, name: "sv_p1", count: 1},
				{dstPort: 68, protocol: 17, name: "sv_p2", count: 0},
				{dstPort: 443, protocol: 6, name: "sv_p2", count: 1},
				{dstPort: 110, protocol: 6, name: "email", count: 1},
				{dstPort: 143, protocol: 6, name: "email", count: 1},
				{dstPort: 993, protocol: 6, name: "email", count: 1},
			},
			map[string]int{
				"sv_p1": 1,
				"sv_p2": 1,
				"email": 3,
			},
			9,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, untagged := stat(strings.NewReader(test.input), test.tags)
			require.Equal(t, test.m, m)
			require.Equal(t, test.untagged, untagged)
		})
	}
}
