package cache

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/seaio-co/util/memory/z"
)

var wait = time.Millisecond * 10

func TestCacheKeyToHash(t *testing.T) {
	keyToHashCount := 0
	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
		KeyToHash: func(key interface{}) (uint64, uint64) {
			keyToHashCount++
			return z.KeyToHash(key)
		},
	})
	require.NoError(t, err)
	if c.Set(1, 1, 1) {
		time.Sleep(wait)
		val, ok := c.Get(1)
		require.True(t, ok)
		require.NotNil(t, val)
		c.Del(1)
	}
	require.Equal(t, 3, keyToHashCount)
}

func TestCacheMaxCost(t *testing.T) {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	key := func() []byte {
		k := make([]byte, 2)
		for i := range k {
			k[i] = charset[rand.Intn(len(charset))]
		}
		return k
	}
	c, err := NewCache(&Config{
		NumCounters: 12960, // 36^2 * 10
		MaxCost:     1e6,   // 1mb
		BufferItems: 64,
		Metrics:     true,
	})
	require.NoError(t, err)
	stop := make(chan struct{}, 8)
	for i := 0; i < 8; i++ {
		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					time.Sleep(time.Millisecond)

					k := key()
					if _, ok := c.Get(k); !ok {
						val := ""
						if rand.Intn(100) < 10 {
							val = "test"
						} else {
							val = strings.Repeat("a", 1000)
						}
						c.Set(key(), val, int64(2+len(val)))
					}
				}
			}
		}()
	}
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second)
		cacheCost := c.Metrics.CostAdded() - c.Metrics.CostEvicted()
		t.Logf("total cache cost: %d\n", cacheCost)
		require.True(t, float64(cacheCost) <= float64(1e6*1.05))
	}
	for i := 0; i < 8; i++ {
		stop <- struct{}{}
	}
}

func TestNewCache(t *testing.T) {
	_, err := NewCache(&Config{
		NumCounters: 0,
	})
	require.Error(t, err)

	_, err = NewCache(&Config{
		NumCounters: 100,
		MaxCost:     0,
	})
	require.Error(t, err)

	_, err = NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 0,
	})
	require.Error(t, err)

	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
		Metrics:     true,
	})
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestNilCache(t *testing.T) {
	var c *Cache
	val, ok := c.Get(1)
	require.False(t, ok)
	require.Nil(t, val)

	require.False(t, c.Set(1, 1, 1))
	c.Del(1)
	c.Clear()
	c.Close()
}

func TestMultipleClose(t *testing.T) {
	var c *Cache
	c.Close()

	var err error
	c, err = NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
		Metrics:     true,
	})
	require.NoError(t, err)
	c.Close()
	c.Close()
}
