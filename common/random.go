package common

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
)

// NewRandom creates a new padded Encoding defined by the given alphabet string.
func NewRandom(alphabet string) *Random {
	r := new(Random)
	diff := 64 - len(alphabet)
	if diff < 0 {
		r.substitute = []byte(alphabet[64:])
		r.substituteLen = len(r.substitute)
		alphabet = alphabet[:64]
	} else {
		r.substitute = []byte(alphabet)
		r.substituteLen = len(r.substitute)
		if diff > 0 {
			alphabet += string(bytes.Repeat([]byte{0x00}, diff))
		}
	}
	r.encoding = base64.NewEncoding(alphabet).WithPadding(base64.NoPadding)
	return r
}

// Random random string creater.
type Random struct {
	encoding      *base64.Encoding
	substitute    []byte
	substituteLen int
}

const urlEncoder = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

var urlRandom = &Random{
	encoding:      base64.URLEncoding,
	substitute:    []byte(urlEncoder),
	substituteLen: len(urlEncoder),
}

// RandomBytes returns securely generated random bytes. It will panic
// if the system's secure random number generator fails to function correctly.
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic(err)
	}
	return b
}
