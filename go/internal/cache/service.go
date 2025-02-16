package cache

import (
	"context"
	"time"

	cacheapiv1 "github.com/ahmedalhulaibi/cache-api/internal/gen/cacheapi/v1"
	"github.com/ahmedalhulaibi/loggy"
)

type cacheService struct {
	logger  *loggy.Logger
	buckets *buckets
	cacheapiv1.UnimplementedCacheServiceServer
}

func NewCacheService(
	logger *loggy.Logger,
) *cacheService {
	return &cacheService{
		logger:  logger,
		buckets: NewCache(),
	}
}

var _ cacheapiv1.CacheServiceServer = (*cacheService)(nil)

func (c *cacheService) Set(ctx context.Context, r *cacheapiv1.SetRequest) (*cacheapiv1.SetResponse, error) {
	c.logger.Infow(ctx, "setting key", "key", r.Key, "bucket", r.Bucket, "value", r.Value)

	ttl := time.Duration(-1)
	evictionPolicy := EvictLRU
	if r.Options != nil {
		ttl = time.Duration(r.Options.TtlSeconds) * time.Second
		evictionPolicy = getEvictionPolicy(r.Options.EvictionPolicy)
	}

	if err := c.buckets.Set(r.Bucket, r.Key, []byte(r.Value), WithTTL(ttl), WithEvictionPolicy(evictionPolicy)); err != nil {
		c.logger.Errorf(ctx, "failed to set key: %v", err)
		return nil, err
	}
	return &cacheapiv1.SetResponse{}, nil
}

func (c *cacheService) Get(ctx context.Context, r *cacheapiv1.GetRequest) (*cacheapiv1.GetResponse, error) {
	record, err := c.buckets.Get(r.Bucket, r.Key)
	if err != nil {
		c.logger.Errorf(ctx, "failed to get key: %v", err)
		return nil, err
	}
	return &cacheapiv1.GetResponse{Value: string(record)}, nil
}

func (c *cacheService) GetStats(ctx context.Context, r *cacheapiv1.GetStatsRequest) (*cacheapiv1.GetStatsResponse, error) {
	s := c.buckets.Stats()
	return &cacheapiv1.GetStatsResponse{
		Hits:      s.Hits,
		Misses:    s.Misses,
		Evictions: s.Evictions,
		Expired:   s.Expired,
	}, nil
}

/*
EvictionPolicy_EVICTION_UNSPECIFIED
EvictionPolicy_EVICTION_LRU
EvictionPolicy_EVICTION_MRU
EvictionPolicy_EVICTION_OLDEST
EvictionPolicy_EVICTION_NEWEST
*/
func getEvictionPolicy(ep cacheapiv1.EvictionPolicy) EvictionPolicy {
	switch ep {
	case cacheapiv1.EvictionPolicy_EVICTION_OLDEST:
		return EvictOldest
	case cacheapiv1.EvictionPolicy_EVICTION_LRU:
		return EvictLRU
	case cacheapiv1.EvictionPolicy_EVICTION_MRU:
		return EvictMRU
	case cacheapiv1.EvictionPolicy_EVICTION_NEWEST:
		return EvictNewest
	default:
		return EvictLRU
	}
}
