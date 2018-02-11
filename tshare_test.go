package tshare

import (
	"bytes"
	"math/rand"
	"testing"
)

type testCase struct {
	Msg    []byte
	Err    error
	Shares [][]byte
}

var testCases []testCase

func init() {
	for i := 0; i <= 4096; i++ {
		msg := make([]byte, i)
		rand.Read(msg)
		shares, err := Split(msg)
		testCases = append(testCases, testCase{
			Msg:    msg,
			Err:    err,
			Shares: shares,
		})
	}
}

func TestSplit(t *testing.T) {
	for i, tc := range testCases {
		if tc.Err != nil {
			t.Errorf("error testCases[%d]", i)
		}
		if len(tc.Shares) != 3 {
			t.Errorf("incorrect number of shares testCases[%d]", i)
		}
		for j := 0; j < 3; j++ {
			if tc.Shares[j][0] != byte(j) {
				t.Errorf("bad tag testCases[%d].Shares[%d]", i, j)
			}
		}
	}
}

func TestJoin(t *testing.T) {
	// Test Join for each permutation of shares
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// Exclude case where only one share is used
			if i == j {
				continue
			}
			testJoin(t, i, j)
		}
	}
}

func testJoin(t *testing.T, i, j int) {
	for n, tc := range testCases {
		a := tc.Shares[i]
		b := tc.Shares[j]
		m, err := Join(a, b)
		if err != nil {
			t.Errorf("Join(s%d, s%d) error testCases[%d]", i, j, n)
		}
		if bytes.Compare(m, tc.Msg) != 0 {
			t.Errorf("Join(s%d, s%d) incorrect testCases[%d]", i, j, n)
		}
	}
}
