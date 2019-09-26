package common

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"github.com/seaio-co/util/stringutil"
	mrand "math/rand"
	"fmt"
)

// NewRandom
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

// Random
type Random struct {
	encoding      *base64.Encoding
	substitute    []byte
	substituteLen int
}

// RandomString
func (r *Random) RandomString(n int) string {
	d := r.encoding.DecodedLen(n)
	buf := make([]byte, n)
	r.encoding.Encode(buf, RandomBytes(d))
	for k, v := range buf {
		if v == 0x00 {
			buf[k] = r.substitute[mrand.Intn(r.substituteLen)]
		}
	}
	return stringutil.BytesToString(buf)
}

const urlEncoder = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

var urlRandom = &Random{
	encoding:      base64.URLEncoding,
	substitute:    []byte(urlEncoder),
	substituteLen: len(urlEncoder),
}

// URLRandomString
func URLRandomString(n int) string {
	return urlRandom.RandomString(n)
}

// RandomBytes
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// GenerateRandnum
func GenerateRandnum() int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(100)

	fmt.Printf("rand is %v\n", randNum)

	return randNum
}

// GenerateRangeNum
func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max - min)
	randNum = randNum + min
	fmt.Printf("rand is %v\n", randNum)
	return randNum
}
