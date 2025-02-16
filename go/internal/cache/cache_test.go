package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCacheImplementationGet(t *testing.T) {
	c := newCache(1)
	defaultOpts, _ := getOptions()

	require.NoError(t, c.Set("user:1", []byte("user1"), defaultOpts))
	record, err := c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Equal(t, []byte("user1"), record)

	// By default, eviction is disabled right now, so this will return an error
	require.Error(t, c.Set("user:2", []byte("user2"), defaultOpts))
	// This record will not be evicted yet since the required at capacity Oldest policy is not applied
	record, err = c.Get("user:1", defaultOpts)
	require.NoError(t, err)
	require.Equal(t, []byte("user1"), record)
}
