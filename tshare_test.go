package tshare

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestJoinBytes(t *testing.T) {
	m := make([]byte, 1024)
	rand.Read(m)
	s, err := SplitBytes(m)
	if err != nil {
		t.Errorf("SplitBytes error: %s", err)
	}
	testJoinBytes(t, m, s[0], s[1])
	testJoinBytes(t, m, s[1], s[0])
	testJoinBytes(t, m, s[0], s[2])
	testJoinBytes(t, m, s[2], s[0])
	testJoinBytes(t, m, s[1], s[2])
	testJoinBytes(t, m, s[2], s[1])
}

func testJoinBytes(t *testing.T, m, a, b []byte) {
	m2, err := JoinBytes(a, b)
	if err != nil {
		t.Errorf("JoinBytes error: %s", err)
	}
	if bytes.Compare(m, m2) != 0 {
		t.Error("JoinBytes incorrect")
	}
}
