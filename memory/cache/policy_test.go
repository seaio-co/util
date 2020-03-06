package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPolicy(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	newPolicy(100, 10)
}

func TestPolicyMetrics(t *testing.T) {
	p := newDefaultPolicy(100, 10)
	p.CollectMetrics(newMetrics())
	require.NotNil(t, p.metrics)
	require.NotNil(t, p.evict.metrics)
}

func TestPolicyProcessItems(t *testing.T) {
	p := newDefaultPolicy(100, 10)
	p.itemsCh <- []uint64{1, 2, 2}
	time.Sleep(wait)
	p.Lock()
	require.Equal(t, int64(2), p.admit.Estimate(2))
	require.Equal(t, int64(1), p.admit.Estimate(1))
	p.Unlock()

	p.stop <- struct{}{}
	p.itemsCh <- []uint64{3, 3, 3}
	time.Sleep(wait)
	p.Lock()
	require.Equal(t, int64(0), p.admit.Estimate(3))
	p.Unlock()
}

func TestPolicyPush(t *testing.T) {
	p := newDefaultPolicy(100, 10)
	require.True(t, p.Push([]uint64{}))

	keepCount := 0
	for i := 0; i < 10; i++ {
		if p.Push([]uint64{1, 2, 3, 4, 5}) {
			keepCount++
		}
	}
	require.NotEqual(t, 0, keepCount)
}

func TestPolicyAdd(t *testing.T) {
	p := newDefaultPolicy(1000, 100)
	if victims, added := p.Add(1, 101); victims != nil || added {
		t.Fatal("can't add an item bigger than entire cache")
	}
	p.Lock()
	p.evict.add(1, 1)
	p.admit.Increment(1)
	p.admit.Increment(2)
	p.admit.Increment(3)
	p.Unlock()

	victims, added := p.Add(1, 1)
	require.Nil(t, victims)
	require.False(t, added)

	victims, added = p.Add(2, 20)
	require.Nil(t, victims)
	require.True(t, added)

	victims, added = p.Add(3, 90)
	require.NotNil(t, victims)
	require.True(t, added)

	victims, added = p.Add(4, 20)
	require.NotNil(t, victims)
	require.False(t, added)
}

func TestPolicyHas(t *testing.T) {
	p := newDefaultPolicy(100, 10)
	p.Add(1, 1)
	require.True(t, p.Has(1))
	require.False(t, p.Has(2))
}
