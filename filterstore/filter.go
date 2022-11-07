package filterstore

import (
	"context"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"github.com/go-redis/redis/v8"
)

type BoolMapPromise interface {
	Result() (map[string]bool, error)
}

type FilterStatusPromise interface {
	Result() (*itemfilterRPC_pb.RedisFilterStatus, error)
}

type Filter interface {
	Create(ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting) error
	Add(ctx context.Context, RefreshTTL bool, elements ...string) error
	CreateAndAdd(ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting, elements ...string) error
	Status(ctx context.Context) (*itemfilterRPC_pb.RedisFilterStatus, error)
	CheckIn(ctx context.Context, candidates ...string) (map[string]bool, error)
	Clean(ctx context.Context) error

	PipeCreateAndAdd(pipe redis.Pipeliner, ctx context.Context, info *itemfilterRPC_pb.RedisFilterSetting, expat string, elements ...string) error
	PipeCheckIn(pipe redis.Pipeliner, ctx context.Context, candidates ...string) (BoolMapPromise, error)
	PipeStatus(pipe redis.Pipeliner, ctx context.Context) FilterStatusPromise
}
