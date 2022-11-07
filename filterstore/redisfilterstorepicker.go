package filterstore

import (
	"context"
	"fmt"
	"time"

	"github.com/Golang-Tools/namespace"
	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"
)

// opr

func (s *RedisFilterStore) PickerKeyTemplate() string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase,
		s.Opts.NamespaceConfig_EntitySourceTypeSubnamespace,
		"{{ EST }}",
		s.Opts.NamespaceConfig_PickerIntervalSubnamespace,
		"{{ PI }}",
		s.Opts.NamespaceConfig_PickerSubnamespace,
		"{{ PID }}"}
	key := ns.ToString(namespace.WithRedisStyle())
	return key
}
func (s *RedisFilterStore) PickerKey(EntitySourceType, PickerInterval, PickerID string) string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase,
		s.Opts.NamespaceConfig_EntitySourceTypeSubnamespace,
		EntitySourceType,
		s.Opts.NamespaceConfig_PickerIntervalSubnamespace,
		PickerInterval,
		s.Opts.NamespaceConfig_PickerSubnamespace,
		PickerID}
	key := ns.ToString(namespace.WithRedisStyle())
	return key
}

func (s *RedisFilterStore) newPicker(ctx context.Context, EntitySourceType, PickerID, TimeRange string, setting *itemfilterRPC_pb.RedisFilterSetting, Content ...string) error {
	key := s.PickerKey(EntitySourceType, TimeRange, PickerID)
	var filter Filter
	if setting.GetBloomfilterSetting() != nil {
		var err error
		fset := setting.GetBloomfilterSetting()
		if fset == nil {
			return ErrFilterSettingNotMatchFilterType
		}
		filter, err = NewBloomFilter(s.conn, key, time.Duration(0))
		if err != nil {
			return err
		}
	} else if setting.GetCuckoofilterSetting() != nil {
		var err error
		fset := setting.GetCuckoofilterSetting()
		if fset == nil {
			return ErrFilterSettingNotMatchFilterType
		}
		filter, err = NewCuckooFilter(s.conn, key, time.Duration(0))
		if err != nil {
			return err
		}
	} else {
		var err error
		fset := setting.GetSetfilterSetting()
		if fset == nil {
			return ErrFilterSettingNotMatchFilterType
		}
		filter, err = NewSetFilter(s.conn, key, time.Duration(0))
		if err != nil {
			return err
		}
	}

	if len(Content) > 0 {
		err := filter.CreateAndAdd(ctx, setting, Content...)
		if err != nil {
			return err
		}
	} else {
		err := filter.Create(ctx, setting)
		if err != nil {
			return err
		}
	}
	at := ExpAtTime(TimeRange)
	if !at.IsZero() {
		s.conn.ExpireAt(ctx, key, at)
	}
	return nil
}

func (s *RedisFilterStore) NewPicker(EntitySourceType, PickerID string, settings map[string]*itemfilterRPC_pb.RedisFilterSetting, Content ...string) error {
	for TimeRange := range settings {
		_, ok := s.Opts.DefaultPickerFilterInfos[TimeRange]
		if !ok {
			return fmt.Errorf("TimeRange %s not Support", TimeRange)
		}
	}
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	for TimeRange, dsetting := range s.Opts.DefaultPickerFilterInfos {
		setting, ok := settings[TimeRange]
		if !ok {
			setting = dsetting
		} else {
			if setting == nil {
				setting = dsetting
			} else {
				if setting.GetBloomfilterSetting() == nil && setting.GetSetfilterSetting() == nil && setting.GetCuckoofilterSetting() == nil {
					setting = dsetting
				}
			}
		}
		err := s.newPicker(ctx, EntitySourceType, PickerID, TimeRange, setting, Content...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *RedisFilterStore) PickerStatus(EntitySourceType, PickerID string) (map[string]*itemfilterRPC_pb.RedisFilterStatus, error) {
	pipe := s.conn.TxPipeline()
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	FilterPromises := map[string]*FilterPromise{}
	for TimeRange := range s.Opts.DefaultPickerFilterInfos {
		key := s.PickerKey(EntitySourceType, TimeRange, PickerID)
		FilterPromises[TimeRange] = findFilterFromKeyPipe(pipe, ctx, key, time.Duration(0), "")
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	pipe = s.conn.TxPipeline()
	resultpromises := map[string]FilterStatusPromise{}
	for TimeRange, filterpromise := range FilterPromises {
		filter, err := filterpromise.Result(s.conn)
		if err != nil {
			return nil, err
		}
		resultpromises[TimeRange] = filter.PipeStatus(pipe, ctx)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	result := map[string]*itemfilterRPC_pb.RedisFilterStatus{}
	for TimeRange, promise := range resultpromises {
		r, err := promise.Result()
		if err != nil {
			return nil, err
		}
		result[TimeRange] = r
	}
	return result, nil
}
