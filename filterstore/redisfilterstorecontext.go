package filterstore

import (
	"time"

	"github.com/Golang-Tools/namespace"
	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"
)

// context
func (s *RedisFilterStore) ContextKeyTemplate() string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.NamespaceConfig_EntitySourceTypeSubnamespace, "{{ EST }}", s.Opts.ContextConfig_KeySubnamespace, "{{ CID }}"}
	key := ns.ToString(namespace.WithRedisStyle())
	return key
}

func (s *RedisFilterStore) ContextKey(EntitySourceType, ContextID string) string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.NamespaceConfig_EntitySourceTypeSubnamespace, EntitySourceType, s.Opts.ContextConfig_KeySubnamespace, ContextID}
	key := ns.ToString(namespace.WithRedisStyle())
	return key
}

func (s *RedisFilterStore) NewContext(EntitySourceType, ContextID string, setting *itemfilterRPC_pb.RedisFilterSetting, Content ...string) error {
	key := s.ContextKey(EntitySourceType, ContextID)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	var filter Filter
	useDefault := false
	//当不填入设置时使用默认设置
	if setting.GetBloomfilterSetting() == nil && setting.GetSetfilterSetting() == nil && setting.GetCuckoofilterSetting() == nil {
		useDefault = true
	}
	if useDefault {
		if s.Opts.DefaultContextFilterInfo.GetBloomfilterSetting() != nil {
			var err error
			fset := s.Opts.DefaultContextFilterInfo.GetBloomfilterSetting()
			TTLSeconds := fset.MaxTTLSeconds
			if TTLSeconds <= 0 {
				filter, err = NewBloomFilter(s.conn, key, time.Duration(0))
			} else {
				filter, err = NewBloomFilter(s.conn, key, time.Duration(TTLSeconds)*time.Second)
			}
			if err != nil {
				return err
			}
		} else if s.Opts.DefaultContextFilterInfo.GetCuckoofilterSetting() != nil {
			var err error
			fset := s.Opts.DefaultContextFilterInfo.GetCuckoofilterSetting()
			TTLSeconds := fset.MaxTTLSeconds
			if TTLSeconds <= 0 {
				filter, err = NewCuckooFilter(s.conn, key, time.Duration(0))
			} else {
				filter, err = NewCuckooFilter(s.conn, key, time.Duration(TTLSeconds)*time.Second)
			}
			if err != nil {
				return err
			}
		} else {
			var err error
			fset := s.Opts.DefaultContextFilterInfo.GetSetfilterSetting()
			TTLSeconds := fset.MaxTTLSeconds
			if TTLSeconds <= 0 {
				filter, err = NewSetFilter(s.conn, key, time.Duration(0))
			} else {
				filter, err = NewSetFilter(s.conn, key, time.Duration(TTLSeconds)*time.Second)
			}
			if err != nil {
				return err
			}
		}
	} else {
		if setting.GetBloomfilterSetting() != nil {
			var err error
			fset := setting.GetBloomfilterSetting()
			if fset == nil {
				return ErrFilterSettingNotMatchFilterType
			}
			TTLSeconds := fset.MaxTTLSeconds
			if TTLSeconds <= 0 {
				filter, err = NewBloomFilter(s.conn, key, time.Duration(0))
			} else {
				filter, err = NewBloomFilter(s.conn, key, time.Duration(TTLSeconds)*time.Second)
			}
			if err != nil {
				return err
			}
		} else if setting.GetCuckoofilterSetting() != nil {
			var err error
			fset := setting.GetCuckoofilterSetting()
			if fset == nil {
				return ErrFilterSettingNotMatchFilterType
			}
			TTLSeconds := fset.MaxTTLSeconds
			if TTLSeconds <= 0 {
				filter, err = NewCuckooFilter(s.conn, key, time.Duration(0))
			} else {
				filter, err = NewCuckooFilter(s.conn, key, time.Duration(TTLSeconds)*time.Second)
			}
			if err != nil {
				return err
			}
		} else {
			var err error
			fset := setting.GetSetfilterSetting()
			if fset == nil {
				return ErrFilterSettingNotMatchFilterType
			}
			TTLSeconds := fset.MaxTTLSeconds
			if TTLSeconds <= 0 {
				filter, err = NewSetFilter(s.conn, key, time.Duration(0))
			} else {
				filter, err = NewSetFilter(s.conn, key, time.Duration(TTLSeconds)*time.Second)
			}
			if err != nil {
				return err
			}
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
	return nil
}
func (s *RedisFilterStore) ContextStatus(EntitySourceType, ContextID string) (*itemfilterRPC_pb.RedisFilterStatus, error) {
	key := s.ContextKey(EntitySourceType, ContextID)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filter, err := findFilterFromKey(s.conn, ctx, key, time.Duration(0))
	if err != nil {
		return nil, err
	}
	return filter.Status(ctx)
}
