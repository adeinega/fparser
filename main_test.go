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
		{"1", "25,tcp,sv_P1", []tag{{dstPort: 25, protocol: 6, name: "sv_p1"}}, nil},
		{"2", "# comment", []tag{}, nil},
		{"3", "1,2,3,4,5", nil,
			errors.New("record on line 1: wrong number of fields")},
		{"4", "65536,tcp,sv_P1", nil,
			errors.New("invalid port 65536")},
		{"5", "22,bla,sv_P1", nil,
			errors.New("invalid protocol number: bla")},
		{"6",
			`# comment
25,tcp,sv_P1
68,udp,sv_P2
`,
			[]tag{{dstPort: 25, protocol: 6, name: "sv_p1"},
				{dstPort: 68, protocol: 17, name: "sv_p2"}},
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
