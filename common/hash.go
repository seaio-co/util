package common

import (
	"hash"
	"crypto/md5"
	"crypto/sha1"
)

// NewHash
func NewHash(h hash.Hash, space UUID, data []byte, version int) UUID {
	h.Reset()
	h.Write(space[:])
	h.Write(data)
	s := h.Sum(nil)
	var uuid UUID
	copy(uuid[:], s)
	uuid[6] = (uuid[6] & 0x0f) | uint8((version&0xf)<<4)
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return uuid
}

//  NewMD5
func NewMD5(space UUID, data []byte) UUID {
	return NewHash(md5.New(), space, data, 3)
}

// NewSHA1
func NewSHA1(space UUID, data []byte) UUID {
	return NewHash(sha1.New(), space, data, 5)
}