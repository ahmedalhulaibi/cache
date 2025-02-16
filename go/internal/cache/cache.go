package cache

import "time"

/*
### Create an in-memory cache that can hold a maximum of 255 keys and satisfies the interface given below.
System Considerations:
 - If the cache is at capacity and a `Get` method is called, apply an `Oldest` eviction policy.
### Implement a HTTP or gRPC server that can handle the following requests:
- Get a key from the cache
- Set a key in the cache
- Get the cache statistics

System Considerations:
- The cache should be shared across all requests.
- No need to worry about authentication or authorization.

### Bonus Points:
- Unit tests
- Containerization
- Documentation
- Error handling
*/

type Cache interface {
	Set(bucket, key string, value []byte, opts ...Option) error
	Get(bucket, key string, opts ...Option) ([]byte, error)
	Delete(bucket, key string, opts ...Option) error
}

type EvictionPolicy string

const (
	EvictLRU    EvictionPolicy = "LRU"
	EvictMRU    EvictionPolicy = "MRU"
	EvictOldest EvictionPolicy = "Oldest"
	EvictNewest EvictionPolicy = "Newest"
)

type Options struct {
	ttl            time.Duration
	evictionPolicy EvictionPolicy
}

type Option func(*Options) error
