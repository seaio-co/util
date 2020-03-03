package cache

import (
	"math/rand"
	"strings"
	"sync"
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

func TestCacheProcessItems(t *testing.T) {
	m := &sync.Mutex{}
	evicted := make(map[uint64]struct{})
	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
		Cost: func(value interface{}) int64 {
			return int64(value.(int))
		},
		OnEvict: func(key, conflict uint64, value interface{}, cost int64) {
			m.Lock()
			defer m.Unlock()
			evicted[key] = struct{}{}
		},
	})
	require.NoError(t, err)

	var key uint64
	var conflict uint64

	key, conflict = z.KeyToHash(1)
	c.setBuf <- &item{
		flag:     itemNew,
		key:      key,
		conflict: conflict,
		value:    1,
		cost:     0,
	}
	time.Sleep(wait)
	require.True(t, c.policy.Has(1))
	require.Equal(t, int64(1), c.policy.Cost(1))

	key, conflict = z.KeyToHash(1)
	c.setBuf <- &item{
		flag:     itemUpdate,
		key:      key,
		conflict: conflict,
		value:    2,
		cost:     0,
	}
	time.Sleep(wait)
	require.Equal(t, int64(2), c.policy.Cost(1))

	key, conflict = z.KeyToHash(1)
	c.setBuf <- &item{
		flag:     itemDelete,
		key:      key,
		conflict: conflict,
	}
	time.Sleep(wait)
	key, conflict = z.KeyToHash(1)
	val, ok := c.store.Get(key, conflict)
	require.False(t, ok)
	require.Nil(t, val)
	require.False(t, c.policy.Has(1))

	key, conflict = z.KeyToHash(2)
	c.setBuf <- &item{
		flag:     itemNew,
		key:      key,
		conflict: conflict,
		value:    2,
		cost:     3,
	}
	key, conflict = z.KeyToHash(3)
	c.setBuf <- &item{
		flag:     itemNew,
		key:      key,
		conflict: conflict,
		value:    3,
		cost:     3,
	}
	key, conflict = z.KeyToHash(4)
	c.setBuf <- &item{
		flag:     itemNew,
		key:      key,
		conflict: conflict,
		value:    3,
		cost:     3,
	}
	key, conflict = z.KeyToHash(5)
	c.setBuf <- &item{
		flag:     itemNew,
		key:      key,
		conflict: conflict,
		value:    3,
		cost:     5,
	}
	time.Sleep(wait)
	m.Lock()
	require.NotEqual(t, 0, len(evicted))
	m.Unlock()

	defer func() {
		require.NotNil(t, recover())
	}()
	c.Close()
	c.setBuf <- &item{flag: itemNew}
}

func TestCacheGet(t *testing.T) {
	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
		Metrics:     true,
	})
	require.NoError(t, err)

	key, conflict := z.KeyToHash(1)
	i := item{
		key:      key,
		conflict: conflict,
		value:    1,
	}
	c.store.Set(&i)
	val, ok := c.Get(1)
	require.True(t, ok)
	require.NotNil(t, val)

	val, ok = c.Get(2)
	require.False(t, ok)
	require.Nil(t, val)

	// 0.5 and not 1.0 because we tried Getting each item twice
	require.Equal(t, 0.5, c.Metrics.Ratio())

	c = nil
	val, ok = c.Get(0)
	require.False(t, ok)
	require.Nil(t, val)
}

// retrySet calls SetWithTTL until the item is accepted by the cache.
func retrySet(t *testing.T, c *Cache, key, value int, cost int64, ttl time.Duration) {
	for {
		if set := c.SetWithTTL(key, value, cost, ttl); !set {
			time.Sleep(wait)
			continue
		}

		time.Sleep(wait)
		val, ok := c.Get(key)
		require.True(t, ok)
		require.NotNil(t, val)
		require.Equal(t, value, val.(int))
		return
	}
}

func TestCacheSet(t *testing.T) {
	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
		Metrics:     true,
	})
	require.NoError(t, err)

	retrySet(t, c, 1, 1, 1, 0)

	c.Set(1, 2, 2)
	val, ok := c.store.Get(z.KeyToHash(1))
	require.True(t, ok)
	require.Equal(t, 2, val.(int))

	c.stop <- struct{}{}
	for i := 0; i < setBufSize; i++ {
		key, conflict := z.KeyToHash(1)
		c.setBuf <- &item{
			flag:     itemUpdate,
			key:      key,
			conflict: conflict,
			value:    1,
			cost:     1,
		}
	}
	require.False(t, c.Set(2, 2, 1))
	require.Equal(t, uint64(1), c.Metrics.SetsDropped())
	close(c.setBuf)
	close(c.stop)

	c = nil
	require.False(t, c.Set(1, 1, 1))
}

func TestRecacheWithTTL(t *testing.T) {
	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
		Metrics:     true,
	})

	require.NoError(t, err)

	// Set initial value for key = 1
	insert := c.SetWithTTL(1, 1, 1, 5*time.Second)
	require.True(t, insert)
	time.Sleep(2 * time.Second)

	// Get value from cache for key = 1
	val, ok := c.Get(1)
	require.True(t, ok)
	require.NotNil(t, val)
	require.Equal(t, 1, val)

	// Wait for expiration
	time.Sleep(5 * time.Second)

	// The cached value for key = 1 should be gone
	val, ok = c.Get(1)
	require.False(t, ok)
	require.Nil(t, val)

	// Set new value for key = 1
	insert = c.SetWithTTL(1, 2, 1, 5*time.Second)
	require.True(t, insert)
	time.Sleep(2 * time.Second)

	// Get value from cache for key = 1
	val, ok = c.Get(1)
	require.True(t, ok)
	require.NotNil(t, val)
	require.Equal(t, 2, val)
}

func TestCacheSetWithTTL(t *testing.T) {
	m := &sync.Mutex{}
	evicted := make(map[uint64]struct{})
	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
		Metrics:     true,
		OnEvict: func(key, conflict uint64, value interface{}, cost int64) {
			m.Lock()
			defer m.Unlock()
			evicted[key] = struct{}{}
		},
	})
	require.NoError(t, err)

	retrySet(t, c, 1, 1, 1, time.Second)

	// Sleep to make sure the item has expired after execution resumes.
	time.Sleep(2 * time.Second)
	val, ok := c.Get(1)
	require.False(t, ok)
	require.Nil(t, val)

	// Sleep to ensure that the bucket where the item was stored has been cleared
	// from the expiraton map.
	time.Sleep(5 * time.Second)
	m.Lock()
	require.Equal(t, 1, len(evicted))
	_, ok = evicted[1]
	require.True(t, ok)
	m.Unlock()

	// Verify that expiration times are overwritten.
	retrySet(t, c, 2, 1, 1, time.Second)
	retrySet(t, c, 2, 2, 1, 100*time.Second)
	time.Sleep(3 * time.Second)
	val, ok = c.Get(2)
	require.True(t, ok)
	require.Equal(t, 2, val.(int))

	// Verify that entries with no expiration are overwritten.
	retrySet(t, c, 3, 1, 1, 0)
	retrySet(t, c, 3, 2, 1, time.Second)
	time.Sleep(3 * time.Second)
	val, ok = c.Get(3)
	require.False(t, ok)
	require.Nil(t, val)
}

func TestCacheDel(t *testing.T) {
	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
	})
	require.NoError(t, err)

	c.Set(1, 1, 1)
	c.Del(1)
	val, ok := c.Get(1)
	require.False(t, ok)
	require.Nil(t, val)

	c = nil
	defer func() {
		require.Nil(t, recover())
	}()
	c.Del(1)
}

func TestCacheDelWithTTL(t *testing.T) {
	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
	})
	require.NoError(t, err)
	retrySet(t, c, 3, 1, 1, 10*time.Second)
	time.Sleep(1 * time.Second)
	// Delete the item
	c.Del(3)
	// Ensure the key is deleted.
	val, ok := c.Get(3)
	require.False(t, ok)
	require.Nil(t, val)
}

func TestCacheClear(t *testing.T) {
	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
		Metrics:     true,
	})
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		c.Set(i, i, 1)
	}
	time.Sleep(wait)
	require.Equal(t, uint64(10), c.Metrics.KeysAdded())

	c.Clear()
	require.Equal(t, uint64(0), c.Metrics.KeysAdded())

	for i := 0; i < 10; i++ {
		val, ok := c.Get(i)
		require.False(t, ok)
		require.Nil(t, val)
	}
}

func TestCacheMetrics(t *testing.T) {
	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
		Metrics:     true,
	})
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		c.Set(i, i, 1)
	}
	time.Sleep(wait)
	m := c.Metrics
	require.Equal(t, uint64(10), m.KeysAdded())
}

func TestMetrics(t *testing.T) {
	newMetrics()
}

func TestNilMetrics(t *testing.T) {
	var m *Metrics
	for _, f := range []func() uint64{
		m.Hits,
		m.Misses,
		m.KeysAdded,
		m.KeysEvicted,
		m.CostEvicted,
		m.SetsDropped,
		m.SetsRejected,
		m.GetsDropped,
		m.GetsKept,
	} {
		require.Equal(t, uint64(0), f())
	}
}

func TestMetricsAddGet(t *testing.T) {
	m := newMetrics()
	m.add(hit, 1, 1)
	m.add(hit, 2, 2)
	m.add(hit, 3, 3)
	require.Equal(t, uint64(6), m.Hits())

	m = nil
	m.add(hit, 1, 1)
	require.Equal(t, uint64(0), m.Hits())
}

func TestMetricsRatio(t *testing.T) {
	m := newMetrics()
	require.Equal(t, float64(0), m.Ratio())

	m.add(hit, 1, 1)
	m.add(hit, 2, 2)
	m.add(miss, 1, 1)
	m.add(miss, 2, 2)
	require.Equal(t, 0.5, m.Ratio())

	m = nil
	require.Equal(t, float64(0), m.Ratio())
}

func TestMetricsString(t *testing.T) {
	m := newMetrics()
	m.add(hit, 1, 1)
	m.add(miss, 1, 1)
	m.add(keyAdd, 1, 1)
	m.add(keyUpdate, 1, 1)
	m.add(keyEvict, 1, 1)
	m.add(costAdd, 1, 1)
	m.add(costEvict, 1, 1)
	m.add(dropSets, 1, 1)
	m.add(rejectSets, 1, 1)
	m.add(dropGets, 1, 1)
	m.add(keepGets, 1, 1)
	require.Equal(t, uint64(1), m.Hits())
	require.Equal(t, uint64(1), m.Misses())
	require.Equal(t, 0.5, m.Ratio())
	require.Equal(t, uint64(1), m.KeysAdded())
	require.Equal(t, uint64(1), m.KeysUpdated())
	require.Equal(t, uint64(1), m.KeysEvicted())
	require.Equal(t, uint64(1), m.CostAdded())
	require.Equal(t, uint64(1), m.CostEvicted())
	require.Equal(t, uint64(1), m.SetsDropped())
	require.Equal(t, uint64(1), m.SetsRejected())
	require.Equal(t, uint64(1), m.GetsDropped())
	require.Equal(t, uint64(1), m.GetsKept())

	require.NotEqual(t, 0, len(m.String()))

	m = nil
	require.Equal(t, 0, len(m.String()))

	require.Equal(t, "unidentified", stringFor(doNotUse))
}

func TestCacheMetricsClear(t *testing.T) {
	c, err := NewCache(&Config{
		NumCounters: 100,
		MaxCost:     10,
		BufferItems: 64,
		Metrics:     true,
	})
	require.NoError(t, err)

	c.Set(1, 1, 1)
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				c.Get(1)
			}
		}
	}()
	time.Sleep(wait)
	c.Clear()
	stop <- struct{}{}
	c.Metrics = nil
	c.Metrics.Clear()
}

func init() {
	// Set bucketSizeSecs to 1 to avoid waiting too much during the tests.
	bucketDurationSecs = 1
}
