package filterstore

import (
	"context"
	"time"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"github.com/go-redis/redis/v8"
)

type SetFilter struct {
	client redis.UniversalClient
	key    string
	setttl time.Duration
}

func NewSetFilter(client redis.UniversalClient, key string, setttl time.Duration) (*SetFilter, error) {
	c := new(SetFilter)
	c.client = client
	c.key = key
	c.setttl = setttl
	return c, nil
}

func (s *SetFilter) Create(ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting) error {
	return nil
}
func (s *SetFilter) Add(ctx context.Context, RefreshTTL bool, elements ...string) error {
	if s.setttl > 0 && RefreshTTL {
		pipe := s.client.TxPipeline()
		pipe.SAdd(ctx, s.key, elements)
		pipe.Expire(ctx, s.key, s.setttl)
		_, err := pipe.Exec(ctx)
		if err != nil {
			return err
		}
	}
	_, err := s.client.SAdd(ctx, s.key, elements).Result()
	if err != nil {
		return err
	}
	return nil
}

func (s *SetFilter) CreateAndAdd(ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting, elements ...string) error {
	return s.Add(ctx, true, elements...)
}

func (s *SetFilter) Status(ctx context.Context) (*itemfilterRPC_pb.RedisFilterStatus, error) {
	pipe := s.client.TxPipeline()
	scardcmd := pipe.SCard(ctx, s.key)
	ttlcmd := pipe.TTL(ctx, s.key)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	num, err := scardcmd.Result()
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
	if num < 1 {
		return &itemfilterRPC_pb.RedisFilterStatus{
			Exists: false,
		}, nil
	}
	return &itemfilterRPC_pb.RedisFilterStatus{
		Exists:     true,
		Size:       num,
		TTLSeconds: int64(ttl.Seconds()),
	}, nil
}
func (s *SetFilter) CheckIn(ctx context.Context, candidates ...string) (map[string]bool, error) {
	members := []interface{}{}
	for _, c := range candidates {
		members = append(members, c)
	}
	res, err := s.client.SMIsMember(ctx, s.key, members...).Result()
	if err != nil {
		return nil, err
	}
	result := map[string]bool{}
	for i, c := range candidates {
		result[c] = res[i]
	}
	return result, nil
}

func (s *SetFilter) Clean(ctx context.Context) error {
	_, err := s.client.Del(ctx, s.key).Result()
	return err
}

func (s *SetFilter) PipeCreateAndAdd(pipe redis.Pipeliner, ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting, expat string, elements ...string) error {
	pipe.SAdd(ctx, s.key, elements)
	if expat != "" {
		at := ExpAtTime(expat)
		if !at.IsZero() {
			pipe.ExpireAt(ctx, s.key, at)
		}
	} else {
		if s.setttl > 0 {
			pipe.Expire(ctx, s.key, s.setttl)
		}
	}
	return nil
}

type SetBoolMapPlacehold struct {
	Cmd   *redis.BoolSliceCmd
	Items []string
}

func (r *SetBoolMapPlacehold) Result() (map[string]bool, error) {
	res, err := r.Cmd.Result()
	if err != nil {
		return nil, err
	}
	result := map[string]bool{}
	for i, c := range r.Items {
		result[c] = res[i]
	}
	return result, nil
}
func (s *SetFilter) PipeCheckIn(pipe redis.Pipeliner, ctx context.Context, candidates ...string) (BoolMapPromise, error) {

	members := []interface{}{}
	for _, c := range candidates {
		members = append(members, c)
	}
	cmd := pipe.SMIsMember(ctx, s.key, members...)

	result := SetBoolMapPlacehold{
		Cmd:   cmd,
		Items: candidates,
	}
	return &result, nil
}

type SetFilterStatusPromise struct {
	scardcmd *redis.IntCmd
	ttlcmd   *redis.DurationCmd
}

func (p *SetFilterStatusPromise) Result() (*itemfilterRPC_pb.RedisFilterStatus, error) {
	num, err := p.scardcmd.Result()
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
	if num < 1 {
		return &itemfilterRPC_pb.RedisFilterStatus{
			Exists: false,
		}, nil
	}
	return &itemfilterRPC_pb.RedisFilterStatus{
		Exists:     true,
		Size:       num,
		TTLSeconds: int64(ttl.Seconds()),
	}, nil
}

func (s *SetFilter) PipeStatus(pipe redis.Pipeliner, ctx context.Context) FilterStatusPromise {
	scardcmd := pipe.SCard(ctx, s.key)
	ttlcmd := pipe.TTL(ctx, s.key)
	return &SetFilterStatusPromise{
		scardcmd: scardcmd,
		ttlcmd:   ttlcmd,
	}
}
