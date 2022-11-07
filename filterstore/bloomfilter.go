package filterstore

import (
	"context"
	"time"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"github.com/Golang-Tools/optparams"
	bloomfilterhelper "github.com/Golang-Tools/redishelper/v2/exthelper/redisbloomhelper/bloomfilter"
	"github.com/go-redis/redis/v8"
)

type BloomFilter struct {
	filter *bloomfilterhelper.BloomFilter
}

func NewBloomFilter(client redis.UniversalClient, key string, setttl time.Duration) (*BloomFilter, error) {
	c := new(BloomFilter)
	filter, err := bloomfilterhelper.New(client, bloomfilterhelper.WithSpecifiedKey(key), bloomfilterhelper.WithMaxTTL(setttl))
	if err != nil {
		return nil, err
	}
	c.filter = filter
	return c, nil
}

func (s *BloomFilter) Create(ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting) error {
	params := []optparams.Option[bloomfilterhelper.ReserveOpts]{}
	setting := info.GetBloomfilterSetting()
	if setting.NonScaling {
		params = append(params, bloomfilterhelper.ReserveWithNonScaling())
	}
	if setting.Expansion > 0 {
		params = append(params, bloomfilterhelper.ReserveWithExpansion(setting.Expansion))
	}
	return s.filter.Reserve(ctx, setting.Capacity, setting.ErrorRate, params...)
}

func (s *BloomFilter) Add(ctx context.Context, RefreshTTL bool, elements ...string) error {
	params := []optparams.Option[bloomfilterhelper.AddOpts]{}
	if RefreshTTL {
		params = append(params, bloomfilterhelper.AddWithRefreshTTL())
	}
	s.filter.MAddItem(ctx, elements, params...)

	return nil
}

func (s *BloomFilter) CreateAndAdd(ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting, elements ...string) error {
	params := []optparams.Option[bloomfilterhelper.InsertOpts]{}
	setting := info.GetBloomfilterSetting()
	if setting.Capacity > 0 {
		params = append(params, bloomfilterhelper.InsertWithCapacity(setting.Capacity))
	}
	if setting.ErrorRate > 0 {
		params = append(params, bloomfilterhelper.InsertWithErrorRate(setting.ErrorRate))
	}
	if setting.NonScaling {
		params = append(params, bloomfilterhelper.InsertWithNonScaling())
	}
	if setting.Expansion > 0 {
		params = append(params, bloomfilterhelper.InsertWithExpansion(setting.Expansion))
	}
	params = append(params, bloomfilterhelper.InsertWithRefreshTTL())
	_, err := s.filter.Insert(ctx, elements, params...)
	return err
}

func (s *BloomFilter) Status(ctx context.Context) (*itemfilterRPC_pb.RedisFilterStatus, error) {
	pipe := s.filter.Client().TxPipeline()
	key := s.filter.Key()
	infocmd := bloomfilterhelper.InfoPipe(pipe, ctx, key)
	ttlcmd := pipe.TTL(ctx, key)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	info, err := infocmd.Result()
	if err != nil {
		if err == redis.Nil {
			return &itemfilterRPC_pb.RedisFilterStatus{
				Exists: false,
			}, nil
		}
		return nil, err
	}
	ttl, err := ttlcmd.Result()
	if err != nil {
		if err == redis.Nil {
			return &itemfilterRPC_pb.RedisFilterStatus{
				Exists: false,
			}, nil
		}
		return nil, err
	}
	return &itemfilterRPC_pb.RedisFilterStatus{
		Exists:     true,
		Size:       info.NumberOfItemsInserted,
		TTLSeconds: int64(ttl.Seconds()),
	}, nil
}

func (s *BloomFilter) CheckIn(ctx context.Context, candidates ...string) (map[string]bool, error) {
	return s.filter.MExistsItem(ctx, candidates...)
}

func (s *BloomFilter) Clean(ctx context.Context) error {
	return s.filter.Clean(ctx)
}

func (s *BloomFilter) PipeCreateAndAdd(pipe redis.Pipeliner, ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting, expat string, elements ...string) error {
	params := []optparams.Option[bloomfilterhelper.InsertOpts]{}
	setting := info.GetBloomfilterSetting()
	if setting.Capacity > 0 {
		params = append(params, bloomfilterhelper.InsertWithCapacity(setting.Capacity))
	}
	if setting.ErrorRate > 0 {
		params = append(params, bloomfilterhelper.InsertWithErrorRate(setting.ErrorRate))
	}
	if setting.NonScaling {
		params = append(params, bloomfilterhelper.InsertWithNonScaling())
	}
	if setting.Expansion > 0 {
		params = append(params, bloomfilterhelper.InsertWithExpansion(setting.Expansion))
	}
	params = append(params, bloomfilterhelper.InsertWithRefreshTTL())
	_, err := bloomfilterhelper.InsertPipe(pipe, ctx, s.filter.Key(), elements, params...)
	if expat != "" {
		at := ExpAtTime(expat)
		if !at.IsZero() {
			pipe.ExpireAt(ctx, s.filter.Key(), at)
		}
	}
	return err
}
func (s *BloomFilter) PipeCheckIn(pipe redis.Pipeliner, ctx context.Context, candidates ...string) (BoolMapPromise, error) {
	return bloomfilterhelper.MExistsItemPipe(pipe, ctx, s.filter.Key(), candidates...)
}

type BloomFilterStatusPromise struct {
	infocmd *bloomfilterhelper.InfoPlacehold
	ttlcmd  *redis.DurationCmd
}

func (p *BloomFilterStatusPromise) Result() (*itemfilterRPC_pb.RedisFilterStatus, error) {
	info, err := p.infocmd.Result()
	if err != nil {
		if err == redis.Nil {
			return &itemfilterRPC_pb.RedisFilterStatus{
				Exists: false,
			}, nil
		}
		return nil, err
	}
	ttl, err := p.ttlcmd.Result()
	if err != nil {
		if err == redis.Nil {
			return &itemfilterRPC_pb.RedisFilterStatus{
				Exists: false,
			}, nil
		}
		return nil, err
	}
	return &itemfilterRPC_pb.RedisFilterStatus{
		Exists:     true,
		Size:       info.NumberOfItemsInserted,
		TTLSeconds: int64(ttl.Seconds()),
	}, nil
}

func (s *BloomFilter) PipeStatus(pipe redis.Pipeliner, ctx context.Context) FilterStatusPromise {
	key := s.filter.Key()
	infocmd := bloomfilterhelper.InfoPipe(pipe, ctx, key)
	ttlcmd := pipe.TTL(ctx, key)
	return &BloomFilterStatusPromise{
		infocmd: infocmd,
		ttlcmd:  ttlcmd,
	}
}
