// Package tshare implements methods for (2,3) threshold secret sharing for
// splitting secrets into three shares. None of the shares alone give away any
// information about the secret (other than the length) but any combination of
// two shares is able to fully recover the secret.
package tshare

import (
	"crypto/rand"
	"errors"
	"io"
)

const bufferSize = 1024

// Split reads from reader and splits data using (2,3) threshold secret sharing
// algorithm, writing each share to the respective writer.
func Split(r io.Reader, w0 io.Writer, w1 io.Writer, w2 io.Writer) error {
	m := make([]byte, bufferSize)
	s0 := make([]byte, bufferSize)
	s1 := make([]byte, bufferSize)
	s2 := make([]byte, bufferSize)
	for {
		n, _ := r.Read(m)
		if n == 0 {
			break
		}
		err := SplitBytes(m, s0, s1, s2)
		if err != nil {
			return err
		}
		_, err = w0.Write(s0[:n])
		if err != nil {
			return err
		}
		_, err = w1.Write(s1[:n])
		if err != nil {
			return err
		}
		_, err = w2.Write(s2[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

// Join recovers a secret message from two or more shares. Any one of the
// supplied readers can be nil and the secret will still be recovered as long as
// at least two are used. The reconstructed secret is written to w.
func Join(w io.Writer, r0 io.Reader, r1 io.Reader, r2 io.Reader) error {
	var m, s0, s1, s2 []byte
	m = make([]byte, bufferSize)
	if r0 != nil {
		s0 = make([]byte, bufferSize)
	}
	if r1 != nil {
		s1 = make([]byte, bufferSize)
	}
	if r2 != nil {
		s2 = make([]byte, bufferSize)
	}
	for {
		n, err := read(s0, s1, s2, r0, r1, r2)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}
		var t0, t1, t2 []byte
		if s0 != nil {
			t0 = s0[:n]
		}
		if s1 != nil {
			t1 = s1[:n]
		}
		if s2 != nil {
			t2 = s2[:n]
		}
		err = JoinBytes(m, t0, t1, t2)
		if err != nil {
			return err
		}
		_, err = w.Write(m[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

// read from all readers where the corresponding share array isn't nil. Returns
// number of bytes read, ensuring that they match and that at least two of the
// shares are non-nil. Does not return EOF error so consumer should check byte
// count for 0 to end stream.
func read(s0, s1, s2 []byte, r0, r1, r2 io.Reader) (int, error) {
	var n0, n1, n2 int
	if s0 != nil {
		n0, _ = io.ReadFull(r0, s0)
	}
	if s1 != nil {
		n1, _ = io.ReadFull(r1, s1)
	}
	if s2 != nil {
		n2, _ = io.ReadFull(r2, s2)
	}
	// verify share sizes match
	if s0 != nil && s1 != nil && s2 != nil {
		if n0 != n1 && n1 != n2 {
			return 0, errors.New("share sizes must match")
		}
		return n0, nil
	}
	if s0 != nil && s1 != nil {
		if n0 != n1 {
			return 0, errors.New("share sizes must match")
		}
		return n0, nil
	}
	if s0 != nil && s2 != nil {
		if n0 != n2 {
			return 0, errors.New("share sizes must match")
		}
		return n0, nil
	}
	if s1 != nil && s2 != nil {
		if n1 != n2 {
			return 0, errors.New("share sizes must match")
		}
		return n1, nil
	}
	// less than two shares supplied are non-nil
	return 0, errors.New("insufficient combination of shares")
}

// SplitBytes splits m into 3 shares using 2,3 threshold secret sharing
// algorithm defined by:
//   m = secret byte with bits [m7 m6 m5 m4 m3 m2 m1 m0]
//   r = random byte
//   s0 = [ 0  0  0  0 m7 m6 m5 m4] ^ r
//   s1 = [m3 m2 m1 m0  0  0  0  0] ^ r
//   s2 = [m7 m6 m5 m4 m3 m2 m1 m0] ^ r
func SplitBytes(m, s0, s1, s2 []byte) error {
	n := len(m)
	r := make([]byte, n)
	_, err := rand.Read(r)
	if err != nil {
		return err
	}
	if len(s0) < n {
		return errors.New("s0 not large enough")
	}
	if len(s1) < n {
		return errors.New("s1 not large enough")
	}
	if len(s2) < n {
		return errors.New("s2 not large enough")
	}
	for i := 0; i < n; i++ {
		s0[i] = ((m[i] & 0xf0) >> 4) ^ r[i]
		s1[i] = ((m[i] & 0x0f) << 4) ^ r[i]
		s2[i] = m[i] ^ r[i]
	}
	return nil
}

// JoinBytes recovers a secret message from two or more shares. Any one of the
// supplied arrayes can be nil and the secret will still be recovered as long as
// at least two are used. The reconstructed secret is written to m.
func JoinBytes(m, s0, s1, s2 []byte) error {
	if s0 != nil && s1 != nil {
		return joinBytes01(m, s0, s1)
	} else if s0 != nil && s2 != nil {
		return joinBytes02(m, s0, s2)
	} else if s1 != nil && s2 != nil {
		return joinBytes12(m, s1, s2)
	}
	return errors.New("insufficient combination of shares")
}

// joinBytes01, when a = s0 and b = s1
//   c = a ^ b
//   m = [c3 c2 c1 c0 0 0 0 0] | [0 0 0 0 c7 c6 c5 c4]
func joinBytes01(m, a, b []byte) error {
	if len(a) != len(b) {
		return errors.New("share sizes must match")
	}
	if len(m) < len(a) {
		return errors.New("m not large enough")
	}
	for i := 0; i < len(a); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c << 4) & 0xf0) | ((c >> 4) & 0x0f)
	}
	return nil
}

// joinBytes02, when a = s0 and b = s2
//   c = a ^ b
//   m = [0 0 0 0 c7 c6 c5 c4] ^ c
func joinBytes02(m, a, b []byte) error {
	if len(a) != len(b) {
		return errors.New("share sizes must match")
	}
	if len(m) < len(a) {
		return errors.New("m not large enough")
	}
	for i := 0; i < len(a); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c & 0xf0) >> 4) ^ c
	}
	return nil
}

// joinBytes12, when a = s1 and b = s2
//   c = a ^ b
//   m = [c3 c2 c1 c0 0 0 0 0] ^ c
func joinBytes12(m, a, b []byte) error {
	if len(a) != len(b) {
		return errors.New("share sizes must match")
	}
	if len(m) < len(a) {
		return errors.New("m not large enough")
	}
	for i := 0; i < len(a); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c & 0x0f) << 4) ^ c
	}
	return nil
}
