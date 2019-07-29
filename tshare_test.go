package tshare

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestJoinBytes(t *testing.T) {
	m := make([]byte, 1024)
	rand.Read(m)
	s0 := make([]byte, len(m))
	s1 := make([]byte, len(m))
	s2 := make([]byte, len(m))
	err := SplitBytes(m, s0, s1, s2)
	if err != nil {
		t.Errorf("SplitBytes error: %s", err)
	}
	testJoinBytes(t, m, s0, s1, nil)
	testJoinBytes(t, m, s0, nil, s2)
	testJoinBytes(t, m, nil, s1, s2)
	testJoinBytes(t, m, s0, s1, s2)
}

func testJoinBytes(t *testing.T, m, s0, s1, s2 []byte) {
	m2 := make([]byte, len(m))
	err := JoinBytes(m2, s0, s1, s2)
	if err != nil {
		t.Errorf("JoinBytes error: %s", err)
	}
	if bytes.Compare(m, m2) != 0 {
		t.Error("JoinBytes incorrect")
	}
}

func TestJoin01(t *testing.T) {
	msg := make([]byte, 1024)
	rand.Read(msg)
	var mi, mo, s0, s1, s2 bytes.Buffer
	mi.Write(msg)
	err := Split(&mi, &s0, &s1, &s2)
	if err != nil {
		t.Errorf("Split error: %s", err)
	}
	err = Join(&mo, &s0, &s1, nil)
	if err != nil {
		t.Errorf("Join error: %s", err)
	}
	if bytes.Compare(msg, mo.Bytes()) != 0 {
		t.Error("Join incorrect")
	}
}

func TestJoin02(t *testing.T) {
	msg := make([]byte, 1024)
	rand.Read(msg)
	var mi, mo, s0, s1, s2 bytes.Buffer
	mi.Write(msg)
	err := Split(&mi, &s0, &s1, &s2)
	if err != nil {
		t.Errorf("Split error: %s", err)
	}
	err = Join(&mo, &s0, nil, &s2)
	if err != nil {
		t.Errorf("Join error: %s", err)
	}
	if bytes.Compare(msg, mo.Bytes()) != 0 {
		t.Error("Join incorrect")
	}
}

func TestJoin12(t *testing.T) {
	msg := make([]byte, 1024)
	rand.Read(msg)
	var mi, mo, s0, s1, s2 bytes.Buffer
	mi.Write(msg)
	err := Split(&mi, &s0, &s1, &s2)
	if err != nil {
		t.Errorf("Split error: %s", err)
	}
	err = Join(&mo, nil, &s1, &s2)
	if err != nil {
		t.Errorf("Join error: %s", err)
	}
	if bytes.Compare(msg, mo.Bytes()) != 0 {
		t.Error("Join incorrect")
	}
}

func TestJoin012(t *testing.T) {
	msg := make([]byte, 1024)
	rand.Read(msg)
	var mi, mo, s0, s1, s2 bytes.Buffer
	mi.Write(msg)
	err := Split(&mi, &s0, &s1, &s2)
	if err != nil {
		t.Errorf("Split error: %s", err)
	}
	err = Join(&mo, &s0, &s1, &s2)
	if err != nil {
		t.Errorf("Join error: %s", err)
	}
	if bytes.Compare(msg, mo.Bytes()) != 0 {
		t.Error("Join incorrect")
	}
}
