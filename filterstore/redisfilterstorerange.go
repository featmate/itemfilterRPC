package filterstore

import (
	"time"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"github.com/Golang-Tools/namespace"
)

// range
func (s *RedisFilterStore) RangeKey(rangeid string) string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.RangeConfig_KeySubnamespace}
	key := ns.FullName(rangeid, namespace.WithRedisStyle())
	return key
}

func (s *RedisFilterStore) NewRange(rangeid string, setting *itemfilterRPC_pb.RedisFilterSetting, Content ...string) error {
	key := s.RangeKey(rangeid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	var filter Filter
	useDefault := false
	//当不填入设置时使用默认设置
	if setting.GetBloomfilterSetting() == nil && setting.GetSetfilterSetting() == nil && setting.GetCuckoofilterSetting() == nil {
		useDefault = true
	}
	if useDefault {
		if s.Opts.DefaultRangeFilterInfo.GetBloomfilterSetting() != nil {
			var err error
			fset := s.Opts.DefaultRangeFilterInfo.GetBloomfilterSetting()
			TTLSeconds := fset.MaxTTLSeconds
			if TTLSeconds <= 0 {
				filter, err = NewBloomFilter(s.conn, key, time.Duration(0))
			} else {
				filter, err = NewBloomFilter(s.conn, key, time.Duration(TTLSeconds)*time.Second)
			}
			if err != nil {
				return err
			}
		} else if s.Opts.DefaultRangeFilterInfo.GetCuckoofilterSetting() != nil {
			var err error
			fset := s.Opts.DefaultRangeFilterInfo.GetCuckoofilterSetting()
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
			fset := s.Opts.DefaultRangeFilterInfo.GetSetfilterSetting()
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
func (s *RedisFilterStore) RangeStatus(rangeid string) (*itemfilterRPC_pb.RedisFilterStatus, error) {
	key := s.RangeKey(rangeid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filter, err := findFilterFromKey(s.conn, ctx, key, time.Duration(0))
	if err != nil {
		return nil, err
	}
	return filter.Status(ctx)
}
func (s *RedisFilterStore) DeleteRange(rangeid string) error {
	key := s.RangeKey(rangeid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filter, err := findFilterFromKey(s.conn, ctx, key, time.Duration(0))
	if err != nil {
		return err
	}
	return filter.Clean(ctx)
}
func (s *RedisFilterStore) AddRangeSource(rangeid string, Content ...string) error {
	key := s.RangeKey(rangeid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filter, err := findFilterFromKey(s.conn, ctx, key, time.Duration(0))
	if err != nil {
		return err
	}
	return filter.Add(ctx, false, Content...)
}
func (s *RedisFilterStore) CheckSourceInRange(rangeid string, Candidates ...string) (CandidateStatus map[string]bool, err error) {
	key := s.RangeKey(rangeid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filter, err := findFilterFromKey(s.conn, ctx, key, time.Duration(0))
	if err != nil {
		return nil, err
	}
	return filter.CheckIn(ctx, Candidates...)
}
func (s *RedisFilterStore) FilterSourceNotInRange(rangeid string, Need, chunkSize int32, Candidates ...string) (Survivor []string, err error) {
	key := s.RangeKey(rangeid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filter, err := findFilterFromKey(s.conn, ctx, key, time.Duration(0))
	if err != nil {
		return nil, err
	}
	if Need > 0 {
		if chunkSize < 1 {
			resmap, err := filter.CheckIn(ctx, Candidates...)
			if err != nil {
				return nil, err
			}
			eles := flitout(false, resmap)
			if len(eles) > int(Need) {
				return eles[:Need], nil
			}
			return eles, nil
		} else {
			chunks := chunkBy(Candidates, int(chunkSize))
			result := []string{}
			for _, chunk := range chunks {
				resmap, err := filter.CheckIn(ctx, chunk...)
				if err != nil {
					return nil, err
				}
				res := flitout(false, resmap)
				for _, r := range res {
					if len(result) >= int(Need) {
						return result, nil
					} else {
						result = append(result, r)
					}
				}
			}
			return result, nil
		}

	} else {
		resmap, err := filter.CheckIn(ctx, Candidates...)
		if err != nil {
			return nil, err
		}
		return flitout(false, resmap), nil
	}
}
