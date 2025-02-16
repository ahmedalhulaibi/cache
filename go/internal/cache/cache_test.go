package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBucket(t *testing.T) {
	b := NewCache()
	require.NotNil(t, b)

	require.NoError(t, b.Set("bucket1", "key1", []byte("value1")))
	record, err := b.Get("bucket1", "key1")
	require.NoError(t, err)
	require.Equal(t, []byte("value1"), record)

	for i := 0; i < 1000; i++ {
		require.NoError(t, b.Set("bucket1", fmt.Sprintf("r%d", i), []byte("value1"), WithTTL(0*time.Second), WithEvictionPolicy(EvictLRU)))
		if i%255 == 0 {
			r, err := b.Get("bucket1", fmt.Sprintf("r%d", i-255))
			require.NoError(t, err)
			require.Nil(t, r)
		}
	}

}

func TestCacheImplementationGet(t *testing.T) {
	c := newCache(1)
	defaultOpts, _ := getOptions()
	defaultOpts.evictOnGet = false
	defaultOpts.evictionPolicy = EvictDisabled

	require.NoError(t, c.Set("user:1", []byte("user1"), defaultOpts))
	record, err := c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Equal(t, []byte("user1"), record)

	// By default, eviction is disabled right now, so this will return an error
	require.Error(t, c.Set("user:2", []byte("user2"), defaultOpts))
	// This record will not be evicted since evictOnGet is false
	record, err = c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Equal(t, []byte("user1"), record)
}

func TestLruEviction(t *testing.T) {
	c := newCache(1)
	defaultOpts, _ := getOptions()
	defaultOpts.evictOnGet = false
	defaultOpts.evictionPolicy = EvictLRU

	require.NoError(t, c.Set("user:1", []byte("user1"), defaultOpts))
	record, err := c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Equal(t, []byte("user1"), record)

	require.NoError(t, c.Set("user:2", []byte("user2"), defaultOpts))
	// This record will be evicted since evictOnGet is false
	record, err = c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Nil(t, record)
}

func TestLruEvictionIncreasedCapacity(t *testing.T) {
	c := newCache(2)
	defaultOpts, _ := getOptions()
	defaultOpts.evictOnGet = false
	defaultOpts.evictionPolicy = EvictLRU

	require.NoError(t, c.Set("user:1", []byte("user1"), defaultOpts))
	require.NoError(t, c.Set("user:2", []byte("user2"), defaultOpts))
	record, err := c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Equal(t, []byte("user1"), record)
	require.NoError(t, c.Set("user:3", []byte("user3"), defaultOpts))
	record, err = c.Get("user:2", defaultOpts)
	require.NoError(t, err)
	require.Nil(t, record)
}

func TestMruEviction(t *testing.T) {
	c := newCache(2)
	defaultOpts, _ := getOptions()
	defaultOpts.evictOnGet = false
	defaultOpts.evictionPolicy = EvictMRU

	require.NoError(t, c.Set("user:1", []byte("user1"), defaultOpts))
	require.NoError(t, c.Set("user:2", []byte("user2"), defaultOpts))
	record, err := c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Equal(t, []byte("user1"), record)
	require.NoError(t, c.Set("user:3", []byte("user3"), defaultOpts))
	record, err = c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Nil(t, record)
}

func TestTtlExpiry(t *testing.T) {
	c := newCache(1)
	defaultOpts, _ := getOptions()
	defaultOpts.evictOnGet = false
	defaultOpts.ttl = time.Second

	now, err := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	require.NoError(t, err)
	defaultOpts.clock = func() time.Time {
		return now
	}

	require.NoError(t, c.Set("user:1", []byte("user1"), defaultOpts))
	record, err := c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Equal(t, []byte("user1"), record)

	defaultOpts.clock = func() time.Time {
		return now.Add(2 * time.Second)
	}
	record, err = c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Nil(t, record)

	stats := c.Stats()
	require.Equal(t, uint64(1), stats.Hits)
	require.Equal(t, uint64(1), stats.Misses)
	require.Equal(t, uint64(0), stats.Evictions)
	require.Equal(t, uint64(1), stats.Expired)
}

func TestOldestEviction(t *testing.T) {
	c := newCache(1)
	defaultOpts, _ := getOptions()
	defaultOpts.evictOnGet = false
	defaultOpts.evictionPolicy = EvictOldest

	require.NoError(t, c.Set("user:1", []byte("user1"), defaultOpts))
	record, err := c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Equal(t, []byte("user1"), record)

	require.NoError(t, c.Set("user:2", []byte("user2"), defaultOpts))
	// This record will be evicted since evictOnGet is false
	record, err = c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Nil(t, record)

	stats := c.Stats()
	require.Equal(t, uint64(1), stats.Hits)
	require.Equal(t, uint64(1), stats.Misses)
	require.Equal(t, uint64(1), stats.Evictions)
	require.Equal(t, uint64(0), stats.Expired)
}
