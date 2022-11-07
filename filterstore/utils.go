package filterstore

import (
	"context"
	"time"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"
	"github.com/go-redis/redis/v8"
)

// flitout
// @params in bool true过滤掉为True的,false过滤掉为False
func flitout(in bool, target map[string]bool) []string {
	result := []string{}
	for ele, v := range target {
		if in {
			if !v {
				result = append(result, ele)
			}
		} else {
			if v {
				result = append(result, ele)
			}
		}
	}
	return result
}

func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

func findFilterFromKey(client redis.UniversalClient, ctx context.Context, key string, ttl time.Duration) (Filter, error) {
	typename, err := client.Type(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	switch typename {
	case "set":
		{
			return NewSetFilter(client, key, ttl)
		}
	case "MBbloom--":
		{
			return NewBloomFilter(client, key, ttl)
		}
	case "MBbloomCF":
		{
			return NewCuckooFilter(client, key, ttl)
		}
	case "none":
		{
			return nil, ErrKeyNotExist
		}
	}
	return nil, ErrKeyNotFilter
}

type FilterPromise struct {
	key            string
	ttl            time.Duration
	cmd            *redis.StatusCmd
	Defaultsetting *itemfilterRPC_pb.RedisFilterSetting
	Expat          string
}

func (p *FilterPromise) Result(client redis.UniversalClient) (Filter, error) {
	typename, err := p.cmd.Result()
	if err != nil {
		return nil, err
	}
	ttl := p.ttl
	var defaultfilter Filter
	if p.Defaultsetting != nil {
		if p.Defaultsetting.GetBloomfilterSetting() != nil {
			var err error
			fset := p.Defaultsetting.GetBloomfilterSetting()
			if fset == nil {
				return nil, ErrFilterSettingNotMatchFilterType
			}
			if ttl <= 0 {
				if fset.MaxTTLSeconds <= 0 {
					ttl = time.Duration(0)
				} else {
					ttl = time.Duration(fset.MaxTTLSeconds) * time.Second
				}
			}
			defaultfilter, err = NewBloomFilter(client, p.key, ttl)
			if err != nil {
				return nil, err
			}
		} else if p.Defaultsetting.GetCuckoofilterSetting() != nil {
			var err error
			fset := p.Defaultsetting.GetCuckoofilterSetting()
			if fset == nil {
				return nil, ErrFilterSettingNotMatchFilterType
			}
			if ttl <= 0 {
				if fset.MaxTTLSeconds <= 0 {
					ttl = time.Duration(0)
				} else {
					ttl = time.Duration(fset.MaxTTLSeconds) * time.Second
				}
			}
			defaultfilter, err = NewCuckooFilter(client, p.key, ttl)
			if err != nil {
				return nil, err
			}
		} else {
			var err error
			fset := p.Defaultsetting.GetSetfilterSetting()
			if fset == nil {
				return nil, ErrFilterSettingNotMatchFilterType
			}
			if ttl <= 0 {
				if fset.MaxTTLSeconds <= 0 {
					ttl = time.Duration(0)
				} else {
					ttl = time.Duration(fset.MaxTTLSeconds) * time.Second
				}
			}
			defaultfilter, err = NewSetFilter(client, p.key, ttl)
			if err != nil {
				return nil, err
			}
		}
	}
	switch typename {
	case "set":
		{
			return NewSetFilter(client, p.key, ttl)
		}
	case "MBbloom--":
		{
			return NewBloomFilter(client, p.key, ttl)
		}
	case "MBbloomCF":
		{
			return NewCuckooFilter(client, p.key, ttl)
		}
	case "none":
		{
			if p.Defaultsetting == nil {
				return nil, ErrKeyNotExist
			}
			return defaultfilter, nil
		}
	}
	return nil, ErrKeyNotFilter
}

func findFilterFromKeyPipe(pipe redis.Pipeliner, ctx context.Context, key string, ttl time.Duration, expat string, defaultsetting ...*itemfilterRPC_pb.RedisFilterSetting) *FilterPromise {
	cmd := pipe.Type(ctx, key)
	var setting *itemfilterRPC_pb.RedisFilterSetting
	if len(defaultsetting) > 0 {
		setting = defaultsetting[0]
	}
	return &FilterPromise{
		key:            key,
		ttl:            ttl,
		cmd:            cmd,
		Defaultsetting: setting,
		Expat:          expat,
	}
}
func getFirstDayOfWeek(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, offset)
}

func ExpAtTime(TimeRange string) time.Time {
	now := time.Now().In(time.UTC)
	switch TimeRange {
	case "day":
		{
			return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1).Add(-1 * time.Second)
		}
	case "week":
		{
			return getFirstDayOfWeek(now).AddDate(0, 0, 7).Add(-1 * time.Second)
		}
	case "month":
		{
			return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0).Add(-1 * time.Second)
		}
	case "year":
		{
			return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(1, 0, 1).Add(-1 * time.Second)
		}
	default:
		{
			return time.Time{}
		}
	}
}
