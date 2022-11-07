package filterstore

import (
	"context"
	"time"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"github.com/Golang-Tools/optparams"
	cuckoofilterhelper "github.com/Golang-Tools/redishelper/v2/exthelper/redisbloomhelper/cuckoofilter"
	"github.com/go-redis/redis/v8"
)

type CuckooFilter struct {
	filter *cuckoofilterhelper.Cuckoofilter
}

func NewCuckooFilter(client redis.UniversalClient, key string, setttl time.Duration) (*CuckooFilter, error) {
	c := new(CuckooFilter)
	filter, err := cuckoofilterhelper.New(client, cuckoofilterhelper.WithSpecifiedKey(key), cuckoofilterhelper.WithMaxTTL(setttl))
	if err != nil {
		return nil, err
	}
	c.filter = filter
	return c, nil
}

func (s *CuckooFilter) Create(ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting) error {
	params := []optparams.Option[cuckoofilterhelper.ReserveOpts]{}
	setting := info.GetCuckoofilterSetting()
	if setting.BucketSize > 0 {
		params = append(params, cuckoofilterhelper.ReserveWithBucketSize(setting.BucketSize))
	}
	if setting.Expansion > 0 {
		params = append(params, cuckoofilterhelper.ReserveWithExpansion(setting.Expansion))
	}
	if setting.MaxIterations > 0 {
		params = append(params, cuckoofilterhelper.ReserveWithMaxIterations(setting.MaxIterations))
	}
	return s.filter.Reserve(ctx, setting.Capacity, params...)
}

func (s *CuckooFilter) Add(ctx context.Context, RefreshTTL bool, elements ...string) error {
	params := []optparams.Option[cuckoofilterhelper.AddOpts]{}
	if RefreshTTL {
		params = append(params, cuckoofilterhelper.AddWithRefreshTTL())
	}
	s.filter.MAddItem(ctx, elements, params...)

	return nil
}

func (s *CuckooFilter) CreateAndAdd(ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting, elements ...string) error {
	params := []optparams.Option[cuckoofilterhelper.InsertOpts]{}
	setting := info.GetCuckoofilterSetting()
	if setting.Capacity > 0 {
		params = append(params, cuckoofilterhelper.InsertWithCapacity(setting.Capacity))
	}
	if setting.Expansion > 0 {
		params = append(params, cuckoofilterhelper.InsertWithExpansion(setting.Expansion))
	}
	params = append(params, cuckoofilterhelper.InsertWithRefreshTTL())
	_, err := s.filter.Insert(ctx, elements, params...)
	return err
}

func (s *CuckooFilter) Status(ctx context.Context) (*itemfilterRPC_pb.RedisFilterStatus, error) {
	pipe := s.filter.Client().TxPipeline()
	key := s.filter.Key()
	infocmd := cuckoofilterhelper.InfoPipe(pipe, ctx, key)
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

func (s *CuckooFilter) CheckIn(ctx context.Context, candidates ...string) (map[string]bool, error) {
	return s.filter.MExistsItem(ctx, candidates...)
}
func (s *CuckooFilter) Clean(ctx context.Context) error {
	return s.filter.Clean(ctx)
}

func (s *CuckooFilter) PipeCreateAndAdd(pipe redis.Pipeliner, ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting, expat string, elements ...string) error {
	params := []optparams.Option[cuckoofilterhelper.InsertOpts]{}
	setting := info.GetCuckoofilterSetting()
	if setting.Capacity > 0 {
		params = append(params, cuckoofilterhelper.InsertWithCapacity(setting.Capacity))
	}
	if setting.Expansion > 0 {
		params = append(params, cuckoofilterhelper.InsertWithExpansion(setting.Expansion))
	}
	params = append(params, cuckoofilterhelper.InsertWithRefreshTTL())
	_, err := cuckoofilterhelper.InsertPipe(pipe, ctx, s.filter.Key(), elements, params...)
	if expat != "" {
		at := ExpAtTime(expat)
		if !at.IsZero() {
			pipe.ExpireAt(ctx, s.filter.Key(), at)
		}
	}
	return err
}
func (s *CuckooFilter) PipeCheckIn(pipe redis.Pipeliner, ctx context.Context, candidates ...string) (BoolMapPromise, error) {
	return cuckoofilterhelper.MExistsItemPipe(pipe, ctx, s.filter.Key(), candidates...)
}

type CuckooFilterStatusPromise struct {
	infocmd *cuckoofilterhelper.InfoPlacehold
	ttlcmd  *redis.DurationCmd
}

func (p *CuckooFilterStatusPromise) Result() (*itemfilterRPC_pb.RedisFilterStatus, error) {
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

func (s *CuckooFilter) PipeStatus(pipe redis.Pipeliner, ctx context.Context) FilterStatusPromise {
	key := s.filter.Key()
	infocmd := cuckoofilterhelper.InfoPipe(pipe, ctx, key)
	ttlcmd := pipe.TTL(ctx, key)
	return &CuckooFilterStatusPromise{
		infocmd: infocmd,
		ttlcmd:  ttlcmd,
	}
}
