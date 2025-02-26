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
	Stats() stats
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
	clock          func() time.Time
	// for testing purposes override this behaviour
	evictOnGet bool
}

func getOptions(opts ...Option) (*Options, error) {
	o := &Options{
		ttl:            -1,
		evictionPolicy: EvictLRU,
		clock:          time.Now,
		evictOnGet:     true,
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

// for testing purposes
func WithClock(clock func() time.Time) Option {
	return func(o *Options) error {
		o.clock = clock
		return nil
	}
}

/*
Assumptions:
1. The 255 keys limit is for each bucket
2. Given that a cache is at capacity and a `Get` method is called and the Oldest eviction policy is applied, we will still return the value for the key
*/

var _ Cache = (*buckets)(nil)

type buckets struct {
	buckets map[string]cache
	sync.RWMutex
}

func NewCache() *buckets {
	return &buckets{
		buckets: make(map[string]cache),
	}
}

func (b *buckets) Set(bucket, key string, value []byte, opts ...Option) error {
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

func (b *buckets) Get(bucket, key string, opts ...Option) ([]byte, error) {
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

func (b *buckets) Delete(bucket, key string, opts ...Option) error {
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

func (b *buckets) Stats() stats {
	b.RLock()
	defer b.RUnlock()
	var h, m, e, ex uint64
	for _, c := range b.buckets {
		s := c.Stats()
		h += s.Hits
		m += s.Misses
		e += s.Evictions
		ex += s.Expired
	}
	return stats{
		Hits:      h,
		Misses:    m,
		Evictions: e,
		Expired:   ex,
	}
}

type cache interface {
	Set(key string, value []byte, opts *Options) error
	Get(key string, opts *Options) ([]byte, error)
	Delete(key string, opts *Options) error
	Stats() stats
}

func newCache(capacity int) *cacheImplementation {
	return &cacheImplementation{
		ruList:     list.New(),
		ruIndex:    make(map[string]*list.Element, capacity),
		oldestList: list.New(),
		capacity:   capacity,
	}
}

type record struct {
	key    string
	value  []byte
	expiry *time.Time
}

/*
map[key]->oldestList->ruList->record
*/
type cacheImplementation struct {
	ruList     *list.List // doubly linked list, front is most recently used
	ruIndex    map[string]*list.Element
	oldestList *list.List // doubly linked list, front is oldest
	capacity   int
	stats      stats
	sync.RWMutex
}

type stats struct {
	Hits, Misses, Evictions, Expired uint64
}

func (c *cacheImplementation) Set(key string, value []byte, opts *Options) error {
	c.Lock()
	defer c.Unlock()

	if c.ruList.Len() >= c.capacity {
		if err := c.evict(opts.evictionPolicy); err != nil {
			return err
		}
	}

	var expiry *time.Time = nil
	if opts.ttl > 0 {
		t := opts.clock().Add(opts.ttl)
		expiry = &t
	}

	r := &record{
		key:    key,
		value:  value,
		expiry: expiry,
	}

	oe, ok := c.ruIndex[key]
	if ok {
		elem, ok := oe.Value.(*list.Element)
		if ok {
			c.ruList.MoveToFront(elem)
			// the old record is replaced with the new one, old one will be garbage collected
			elem.Value = r
			return nil
		}
	}

	if c.ruList.Len() == 0 {
		c.ruIndex[key] = c.oldestList.PushFront(c.ruList.PushFront(r))
		return nil
	}

	c.ruIndex[key] = c.oldestList.InsertBefore(
		c.ruList.InsertBefore(r, c.ruList.Front()),
		c.oldestList.Front(),
	)

	return nil
}

func (c *cacheImplementation) Get(key string, opts *Options) ([]byte, error) {
	c.Lock()
	defer c.Unlock()

	elem, ok := c.ruIndex[key]
	if !ok {
		c.stats.Misses++
		return nil, nil
	}

	record := elem.Value.(*list.Element).Value.(*record)

	if record.expiry != nil && opts.clock().After(*record.expiry) {
		c.stats.Misses++
		c.stats.Expired++
		c.remove(elem)
		return nil, nil
	}

	if opts.evictOnGet && c.ruList.Len() >= c.capacity {
		if err := c.evict(EvictOldest); err != nil {
			return nil, err
		}
	} else {
		c.ruList.MoveToFront(elem.Value.(*list.Element))
	}

	c.stats.Hits++
	return record.value, nil
}

func (c *cacheImplementation) Delete(key string, opts *Options) error {
	c.Lock()
	defer c.Unlock()

	elem, ok := c.ruIndex[key]
	if !ok {
		return nil
	}

	c.remove(elem)
	return nil
}

func (c *cacheImplementation) Stats() stats {
	c.RLock()
	defer c.RUnlock()
	return c.stats
}

func (c *cacheImplementation) remove(e *list.Element) {
	c.oldestList.Remove(e)
	r := c.ruList.Remove(e.Value.(*list.Element)).(*record)
	delete(c.ruIndex, r.key)
}

func (c *cacheImplementation) evict(e EvictionPolicy) error {
	elem, err := c.getEvictionCandidate(e)
	if err != nil {
		return err
	}

	if elem == nil {
		return nil
	}

	c.remove(elem)
	c.stats.Evictions++
	return nil
}

func (c *cacheImplementation) getEvictionCandidate(e EvictionPolicy) (*list.Element, error) {
	switch e {
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
	r, ok := c.ruList.Back().Value.(*record)
	if !ok {
		return nil, fmt.Errorf("error getting LRU element")
	}
	return c.ruIndex[r.key], nil
}

func (c *cacheImplementation) getMru() (*list.Element, error) {
	r, ok := c.ruList.Front().Value.(*record)
	if !ok {
		return nil, fmt.Errorf("error getting MRU element")
	}
	return c.ruIndex[r.key], nil
}

func (c *cacheImplementation) getOldest() (*list.Element, error) {
	return c.oldestList.Front(), nil
}

func (c *cacheImplementation) getNewest() (*list.Element, error) {
	return c.oldestList.Back(), nil
}
