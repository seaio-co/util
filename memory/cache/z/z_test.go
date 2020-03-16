package z

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func verifyHashProduct(t *testing.T, wantKey, wantConflict, key, conflict uint64) {
	require.Equal(t, wantKey, key)
	require.Equal(t, wantConflict, conflict)
}

func TestKeyToHash(t *testing.T) {
	var key uint64
	var conflict uint64

	key, conflict = KeyToHash(uint64(1))
	verifyHashProduct(t, 1, 0, key, conflict)

	key, conflict = KeyToHash(1)
	verifyHashProduct(t, 1, 0, key, conflict)

	key, conflict = KeyToHash(int32(2))
	verifyHashProduct(t, 2, 0, key, conflict)

	key, conflict = KeyToHash(int32(-2))
	verifyHashProduct(t, math.MaxUint64-1, 0, key, conflict)

	key, conflict = KeyToHash(int64(-2))
	verifyHashProduct(t, math.MaxUint64-1, 0, key, conflict)

	key, conflict = KeyToHash(uint32(3))
	verifyHashProduct(t, 3, 0, key, conflict)

	key, conflict = KeyToHash(int64(3))
	verifyHashProduct(t, 3, 0, key, conflict)
}
