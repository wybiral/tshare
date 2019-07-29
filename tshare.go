// Package tshare implements methods for (2,3) threshold secret sharing for
// splitting secrets into three shares. None of the shares alone give away any
// information about the secret (other than the length) but any combination of
// two shares is able to fully recover the secret.
package tshare

import (
	"crypto/rand"
	"errors"
)

// ErrSizeMismatch returned when share sizes don't match.
var ErrSizeMismatch = errors.New("share size mismatch")

// ErrInvalidShares returned when invalid combination of tagged shares.
var ErrInvalidShares = errors.New("invalid shares")

// SplitBytes splits m into 3 shares using 2,3 threshold secret sharing
// algorithm defined by:
//   m = secret byte with bits [m7 m6 m5 m4 m3 m2 m1 m0]
//   r = random byte
//   s0 = [ 0  0  0  0 m7 m6 m5 m4] ^ r
//   s1 = [m3 m2 m1 m0  0  0  0  0] ^ r
//   s2 = [m7 m6 m5 m4 m3 m2 m1 m0] ^ r
// The first byte of each share is a tag denoting which share it is.
func SplitBytes(m []byte) ([][]byte, error) {
	n := len(m)
	r := make([]byte, n)
	_, err := rand.Read(r)
	if err != nil {
		return nil, err
	}
	// allocate shares with space for tag byte
	_s0 := make([]byte, n+1)
	_s1 := make([]byte, n+1)
	_s2 := make([]byte, n+1)
	// set tags
	_s0[0] = 0
	_s1[0] = 1
	_s2[0] = 2
	// get actual share slices
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

// JoinBytes recovers secret from any two tagged shares.
func JoinBytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, ErrSizeMismatch
	}
	// shares always have tag byte
	if len(a) < 1 {
		return nil, ErrInvalidShares
	}
	// sort by tags
	if a[0] > b[0] {
		a, b = b, a
	}
	// allocate space for assembled message
	m := make([]byte, len(a)-1)
	if a[0] == 0 && b[0] == 1 {
		joinBytes01(m, a[1:], b[1:])
	} else if a[0] == 0 && b[0] == 2 {
		joinBytes02(m, a[1:], b[1:])
	} else if a[0] == 1 && b[0] == 2 {
		joinBytes12(m, a[1:], b[1:])
	} else {
		return nil, ErrInvalidShares
	}
	return m, nil
}

// joinBytes01, when a = s0 and b = s1
//   c = a ^ b
//   m = [c3 c2 c1 c0 0 0 0 0] | [0 0 0 0 c7 c6 c5 c4]
func joinBytes01(m, a, b []byte) {
	for i := 0; i < len(m); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c << 4) & 0xf0) | ((c >> 4) & 0x0f)
	}
}

// joinBytes02, when a = s0 and b = s2
//   c = a ^ b
//   m = [0 0 0 0 c7 c6 c5 c4] ^ c
func joinBytes02(m, a, b []byte) {
	for i := 0; i < len(m); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c & 0xf0) >> 4) ^ c
	}
}

// joinBytes12, when a = s1 and b = s2
//   c = a ^ b
//   m = [c3 c2 c1 c0 0 0 0 0] ^ c
func joinBytes12(m, a, b []byte) {
	for i := 0; i < len(m); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c & 0x0f) << 4) ^ c
	}
}
