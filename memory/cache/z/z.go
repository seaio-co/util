package z

import (
	"github.com/cespare/xxhash"
)

// TODO: Figure out a way to re-use memhash for the second uint64 hash, we
//       already know that appending bytes isn't reliable for generating a
//       second hash (see Ristretto PR #88).
//
//       We also know that while the Go runtime has a runtime memhash128
//       function, it's not possible to use it to generate [2]uint64 or
//       anything resembling a 128bit hash, even though that's exactly what
//       we need in this situation.
func KeyToHash(key interface{}) (uint64, uint64) {
	if key == nil {
		return 0, 0
	}
	switch k := key.(type) {
	case uint64:
		return k, 0
	case string:
		raw := []byte(k)
		return MemHash(raw), xxhash.Sum64(raw)
	case []byte:
		return MemHash(k), xxhash.Sum64(k)
	case byte:
		return uint64(k), 0
	case int:
		return uint64(k), 0
	case int32:
		return uint64(k), 0
	case uint32:
		return uint64(k), 0
	case int64:
		return uint64(k), 0
	default:
		panic("Key type not supported")
	}
}
