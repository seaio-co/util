package ristretto

import (
	"errors"
	"time"

	"github.com/seaio-co/util/memory/z"
)

const (
	// TODO: find the optimal value for this or make it configurable
	setBufSize = 32 * 1024
)

type onEvictFunc func(uint64, uint64, interface{}, int64)

// Cache is a thread-safe implementation of a hashmap with a TinyLFU admission
// policy and a Sampled LFU eviction policy. You can use the same Cache instance
// from as many goroutines as you want.
type Cache struct {
	// store is the central concurrent hashmap where key-value items are stored.
	store store
	// policy determines what gets let in to the cache and what gets kicked out.
	policy policy
	// getBuf is a custom ring buffer implementation that gets pushed to when
	// keys are read.
	getBuf *ringBuffer
	// setBuf is a buffer allowing us to batch/drop Sets during times of high
	// contention.
	setBuf chan *item
	// onEvict is called for item evictions.
	onEvict onEvictFunc
	// KeyToHash function is used to customize the key hashing algorithm.
	// Each key will be hashed using the provided function. If keyToHash value
	// is not set, the default keyToHash function is used.
	keyToHash func(interface{}) (uint64, uint64)
	// stop is used to stop the processItems goroutine.
	stop chan struct{}
	// cost calculates cost from a value.
	cost func(value interface{}) int64
	// cleanupTicker is used to periodically check for entries whose TTL has passed.
	cleanupTicker *time.Ticker
	// Metrics contains a running log of important statistics like hits, misses,
	// and dropped items.
	Metrics *Metrics
}

// Config is passed to NewCache for creating new Cache instances.
type Config struct {
	// NumCounters determines the number of counters (keys) to keep that hold
	// access frequency information. It's generally a good idea to have more
	// counters than the max cache capacity, as this will improve eviction
	// accuracy and subsequent hit ratios.
	//
	// For example, if you expect your cache to hold 1,000,000 items when full,
	// NumCounters should be 10,000,000 (10x). Each counter takes up 4 bits, so
	// keeping 10,000,000 counters would require 5MB of memory.
	NumCounters int64
	// MaxCost can be considered as the cache capacity, in whatever units you
	// choose to use.
	//
	// For example, if you want the cache to have a max capacity of 100MB, you
	// would set MaxCost to 100,000,000 and pass an item's number of bytes as
	// the `cost` parameter for calls to Set. If new items are accepted, the
	// eviction process will take care of making room for the new item and not
	// overflowing the MaxCost value.
	MaxCost int64
	// BufferItems determines the size of Get buffers.
	//
	// Unless you have a rare use case, using `64` as the BufferItems value
	// results in good performance.
	BufferItems int64
	// Metrics determines whether cache statistics are kept during the cache's
	// lifetime. There *is* some overhead to keeping statistics, so you should
	// only set this flag to true when testing or throughput performance isn't a
	// major factor.
	Metrics bool
	// OnEvict is called for every eviction and passes the hashed key, value,
	// and cost to the function.
	OnEvict func(key, conflict uint64, value interface{}, cost int64)
	// KeyToHash function is used to customize the key hashing algorithm.
	// Each key will be hashed using the provided function. If keyToHash value
	// is not set, the default keyToHash function is used.
	KeyToHash func(key interface{}) (uint64, uint64)
	// Cost evaluates a value and outputs a corresponding cost. This function
	// is ran after Set is called for a new item or an item update with a cost
	// param of 0.
	Cost func(value interface{}) int64
}

type itemFlag byte

const (
	itemNew itemFlag = iota
	itemDelete
	itemUpdate
)

// item is passed to setBuf so items can eventually be added to the cache.
type item struct {
	flag       itemFlag
	key        uint64
	conflict   uint64
	value      interface{}
	cost       int64
	expiration time.Time
}

// NewCache returns a new Cache instance and any configuration errors, if any.
func NewCache(config *Config) (*Cache, error) {
	switch {
	case config.NumCounters == 0:
		return nil, errors.New("NumCounters can't be zero")
	case config.MaxCost == 0:
		return nil, errors.New("MaxCost can't be zero")
	case config.BufferItems == 0:
		return nil, errors.New("BufferItems can't be zero")
	}
	policy := newPolicy(config.NumCounters, config.MaxCost)
	cache := &Cache{
		store:         newStore(),
		policy:        policy,
		getBuf:        newRingBuffer(policy, config.BufferItems),
		setBuf:        make(chan *item, setBufSize),
		onEvict:       config.OnEvict,
		keyToHash:     config.KeyToHash,
		stop:          make(chan struct{}),
		cost:          config.Cost,
		cleanupTicker: time.NewTicker(time.Duration(bucketDurationSecs) * time.Second / 2),
	}
	if cache.keyToHash == nil {
		cache.keyToHash = z.KeyToHash
	}
	if config.Metrics {
		cache.collectMetrics()
	}
	// NOTE: benchmarks seem to show that performance decreases the more
	//       goroutines we have running cache.processItems(), so 1 should
	//       usually be sufficient
	go cache.processItems()
	return cache, nil
}

// Get returns the value (if any) and a boolean representing whether the
// value was found or not. The value can be nil and the boolean can be true at
// the same time.
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	if c == nil || key == nil {
		return nil, false
	}
	keyHash, conflictHash := c.keyToHash(key)
	c.getBuf.Push(keyHash)
	value, ok := c.store.Get(keyHash, conflictHash)
	if ok {
		c.Metrics.add(hit, keyHash, 1)
	} else {
		c.Metrics.add(miss, keyHash, 1)
	}
	return value, ok
}

// Set attempts to add the key-value item to the cache. If it returns false,
// then the Set was dropped and the key-value item isn't added to the cache. If
// it returns true, there's still a chance it could be dropped by the policy if
// its determined that the key-value item isn't worth keeping, but otherwise the
// item will be added and other items will be evicted in order to make room.
//
// To dynamically evaluate the items cost using the Config.Coster function, set
// the cost parameter to 0 and Coster will be ran when needed in order to find
// the items true cost.
func (c *Cache) Set(key, value interface{}, cost int64) bool {
	return c.SetWithTTL(key, value, cost, 0*time.Second)
}

// SetWithTTL works like Set but adds a key-value pair to the cache that will expire
// after the specified TTL (time to live) has passed. A zero value means the value never
// expires, which is identical to calling Set. A negative value is a no-op and the value
// is discarded.
func (c *Cache) SetWithTTL(key, value interface{}, cost int64, ttl time.Duration) bool {
	if c == nil || key == nil {
		return false
	}

	var expiration time.Time
	switch {
	case ttl == 0:
		// No expiration.
		break
	case ttl < 0:
		// Treat this a a no-op.
		return false
	default:
		expiration = time.Now().Add(ttl)
	}

	keyHash, conflictHash := c.keyToHash(key)
	i := &item{
		flag:       itemNew,
		key:        keyHash,
		conflict:   conflictHash,
		value:      value,
		cost:       cost,
		expiration: expiration,
	}
	// cost is eventually updated. The expiration must also be immediately updated
	// to prevent items from being prematurely removed from the map.
	if c.store.Update(i) {
		i.flag = itemUpdate
	}
	// Attempt to send item to policy.
	select {
	case c.setBuf <- i:
		return true
	default:
		c.Metrics.add(dropSets, keyHash, 1)
		return false
	}
}

// Del deletes the key-value item from the cache if it exists.
func (c *Cache) Del(key interface{}) {
	if c == nil || key == nil {
		return
	}
	keyHash, conflictHash := c.keyToHash(key)
	// Delete immediately.
	c.store.Del(keyHash, conflictHash)
	// If we've set an item, it would be applied slightly later.
	// So we must push the same item to `setBuf` with the deletion flag.
	// This ensures that if a set is followed by a delete, it will be
	// applied in the correct order.
	c.setBuf <- &item{
		flag:     itemDelete,
		key:      keyHash,
		conflict: conflictHash,
	}
}

// Close stops all goroutines and closes all channels.
func (c *Cache) Close() {
	if c == nil || c.stop == nil {
		return
	}
	// Block until processItems goroutine is returned.
	c.stop <- struct{}{}
	close(c.stop)
	c.stop = nil
	close(c.setBuf)
	c.policy.Close()
}

// Clear empties the hashmap and zeroes all policy counters. Note that this is
// not an atomic operation (but that shouldn't be a problem as it's assumed that
// Set/Get calls won't be occurring until after this).
func (c *Cache) Clear() {
	if c == nil {
		return
	}
	// Block until processItems goroutine is returned.
	c.stop <- struct{}{}
	// Swap out the setBuf channel.
	c.setBuf = make(chan *item, setBufSize)
	// Clear value hashmap and policy data.
	c.policy.Clear()
	c.store.Clear()
	// Only reset metrics if they're enabled.
	if c.Metrics != nil {
		c.Metrics.Clear()
	}
	// Restart processItems goroutine.
	go c.processItems()
}
