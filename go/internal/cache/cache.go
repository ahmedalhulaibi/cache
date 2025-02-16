package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

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
	EvictDisabled EvictionPolicy = "Disabled"
	EvictLRU      EvictionPolicy = "LRU"
	EvictMRU      EvictionPolicy = "MRU"
	EvictOldest   EvictionPolicy = "Oldest"
	EvictNewest   EvictionPolicy = "Newest"
)

type Options struct {
	ttl            time.Duration
	evictionPolicy EvictionPolicy
}

func getOptions(opts ...Option) (*Options, error) {
	o := &Options{
		ttl:            -1,
		evictionPolicy: EvictDisabled,
	}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

type Option func(*Options) error

func WithTTL(ttl time.Duration) Option {
	return func(o *Options) error {
		o.ttl = ttl
		return nil
	}
}

func WithEvictionPolicy(policy EvictionPolicy) Option {
	return func(o *Options) error {
		o.evictionPolicy = policy
		return nil
	}
}

/*
Assumptions:
1. The 255 keys limit is for each bucket
2. Given that a cache is at capacity and a `Get` method is called and the Oldest eviction policy is applied, we will still return the value for the key
*/

var _ Cache = &bucket{}

type bucket struct {
	buckets map[string]cache
	sync.RWMutex
}

func (b *bucket) Set(bucket, key string, value []byte, opts ...Option) error {
	o, err := getOptions(opts...)
	if err != nil {
		return err
	}

	b.Lock()
	defer b.Unlock()
	if _, ok := b.buckets[bucket]; !ok {
		b.buckets[bucket] = newCache(255)
	}
	return b.buckets[bucket].Set(key, value, o)
}

func (b *bucket) Get(bucket, key string, opts ...Option) ([]byte, error) {
	o, err := getOptions(opts...)
	if err != nil {
		return nil, err
	}

	b.RLock()
	defer b.RUnlock()
	if _, ok := b.buckets[bucket]; !ok {
		return nil, nil
	}
	return b.buckets[bucket].Get(key, o)
}

func (b *bucket) Delete(bucket, key string, opts ...Option) error {
	o, err := getOptions(opts...)
	if err != nil {
		return err
	}

	b.RLock()
	defer b.RUnlock()
	if _, ok := b.buckets[bucket]; !ok {
		return nil
	}
	return b.buckets[bucket].Delete(key, o)
}

type cache interface {
	Set(key string, value []byte, opts *Options) error
	Get(key string, opts *Options) ([]byte, error)
	Delete(key string, opts *Options) error
}

func newCache(capacity int) *cacheImplementation {
	return &cacheImplementation{
		recentlyUsed: list.New(),
		// age:          list.New(),
		keys:     make(map[string]*list.Element, capacity),
		capacity: capacity,
	}
}

type record struct {
	key    string
	value  []byte
	expiry *time.Time
}

type cacheImplementation struct {
	recentlyUsed *list.List // doubly linked list, front is most recently used
	// age          *list.List // doubly linked list, front is oldest
	keys     map[string]*list.Element
	capacity int
	sync.RWMutex
}

func (c *cacheImplementation) Set(key string, value []byte, opts *Options) error {
	c.Lock()
	defer c.Unlock()

	if c.recentlyUsed.Len() == c.capacity {
		if err := c.evict(opts); err != nil {
			return err
		}
	}

	var expiry *time.Time = nil
	if opts.ttl > 0 {
		t := time.Now().Add(opts.ttl)
		expiry = &t
	}

	r := &record{
		key:    key,
		value:  value,
		expiry: expiry,
	}

	elem, ok := c.keys[key]
	if ok {
		c.recentlyUsed.MoveToFront(elem)
		// the old record is replaced with the new one, old one will be garbage collected
		elem.Value = r
		return nil
	}

	newElem := list.Element{Value: r}

	c.recentlyUsed.PushFront(&newElem)
	c.keys[key] = &newElem

	return nil
}

func (c *cacheImplementation) Get(key string, opts *Options) ([]byte, error) {
	c.Lock()
	defer c.Unlock()

	elem, ok := c.keys[key]
	if !ok {
		return nil, nil
	}

	c.recentlyUsed.MoveToFront(elem)
	return elem.Value.(*record).value, nil
}

func (c *cacheImplementation) Delete(key string, opts *Options) error {
	panic("unimplemented")
}

func (c *cacheImplementation) evict(opts *Options) error {
	elem, err := c.getEvictionCandidate(opts)
	if err != nil {
		return err
	}

	if elem == nil {
		return nil
	}

	c.recentlyUsed.Remove(elem)
	delete(c.keys, elem.Value.(*record).key)
	return nil
}

func (c *cacheImplementation) getEvictionCandidate(opts *Options) (*list.Element, error) {
	switch opts.evictionPolicy {
	case EvictLRU:
		return c.getLru()
	case EvictMRU:
		return c.getMru()
	case EvictOldest:
		return c.getOldest()
	case EvictNewest:
		return c.getNewest()
	case EvictDisabled:
		return nil, fmt.Errorf("eviction disabled")
	default:
		return nil, fmt.Errorf("eviction policy not implemented")
	}
}

func (c *cacheImplementation) getLru() (*list.Element, error) {
	return c.recentlyUsed.Back(), nil
}

func (c *cacheImplementation) getMru() (*list.Element, error) {
	return c.recentlyUsed.Front(), nil
}

func (c *cacheImplementation) getOldest() (*list.Element, error) {
	// TODO: implement
	return nil, fmt.Errorf("not implemented")
}

func (c *cacheImplementation) getNewest() (*list.Element, error) {
	// TODO: implement
	return nil, fmt.Errorf("not implemented")
}
