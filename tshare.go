package tshare

import (
	"crypto/rand"
	"errors"
)

// Secrets are split bytewise into three shares using the following rule:
//
// m = secret byte with bits [m7 m6 m5 m4 m3 m2 m1 m0]
// r = random byte
// s0 = [ 0  0  0  0 m7 m6 m5 m4] ^ r
// s1 = [m3 m2 m1 m0  0  0  0  0] ^ r
// s2 = [m7 m6 m5 m4 m3 m2 m1 m0] ^ r
//
// The first byte of each share is a tag to denote which share it is.
func Split(m []byte) ([][]byte, error) {
	n := len(m)
	r := make([]byte, n)
	_, err := rand.Read(r)
	if err != nil {
		return nil, err
	}
	// Allocate shares with space for tag byte
	_s0 := make([]byte, n+1)
	_s1 := make([]byte, n+1)
	_s2 := make([]byte, n+1)
	// Set tags
	_s0[0] = 0
	_s1[0] = 1
	_s2[0] = 2
	// Get actual share slices
	s0 := _s0[1:]
	s1 := _s1[1:]
	s2 := _s2[1:]
	for i := 0; i < n; i++ {
		s0[i] = ((m[i] & 0xf0) >> 4) ^ r[i]
		s1[i] = ((m[i] & 0x0f) << 4) ^ r[i]
		s2[i] = m[i] ^ r[i]
	}
	return [][]byte{_s0, _s1, _s2}, nil
}

// Secrets are joined from any two shares based on the tags used.
func Join(a, b []byte) ([]byte, error) {
	// Sort by tags
	if a[0] > b[0] {
		a, b = b, a
	}
	// Allocate space for assembled message
	m := make([]byte, len(a)-1)
	if a[0] == 0 && b[0] == 1 {
		join01(m, a[1:], b[1:])
	} else if a[0] == 0 && b[0] == 2 {
		join02(m, a[1:], b[1:])
	} else if a[0] == 1 && b[0] == 2 {
		join12(m, a[1:], b[1:])
	} else {
		return nil, errors.New("invalid share combination")
	}
	return m, nil
}

// When a = s0 and b = s1
// c = a ^ b
// m = [c3 c2 c1 c0 0 0 0 0] ^ [0 0 0 0 c7 c6 c5 c4]
func join01(m, a, b []byte) {
	for i := 0; i < len(m); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c << 4) & 0xf0) ^ ((c >> 4) & 0x0f)
	}
}

// When a = s0 and b = s2
// c = a ^ b
// m = [0 0 0 0 c7 c6 c5 c4] ^ c
func join02(m, a, b []byte) {
	for i := 0; i < len(m); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c & 0xf0) >> 4) ^ c
	}
}

// When a = s1 and b = s2
// c = a ^ b
// m = [c3 c2 c1 c0 0 0 0 0] ^ c
func join12(m, a, b []byte) {
	for i := 0; i < len(m); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c & 0x0f) << 4) ^ c
	}
}
