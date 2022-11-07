package filterstore

import (
	"time"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"github.com/Golang-Tools/namespace"
)

// 黑名单
func (s *RedisFilterStore) BlacklistKey(blacklistid string) string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.BlacklistConfig_KeySubnamespace}
	key := ns.FullName(blacklistid, namespace.WithRedisStyle())
	return key
}

func (s *RedisFilterStore) NewBlacklist(blacklistid string, setting *itemfilterRPC_pb.RedisFilterSetting, Content ...string) error {
	key := s.BlacklistKey(blacklistid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	var filter Filter
	useDefault := false
	//当不填入设置时使用默认设置
	if setting.GetSetfilterSetting() == nil {
		useDefault = true
	}
	var err error
	if useDefault {
		fset := s.Opts.DefaultBlacklistFilterInfo
		TTLSeconds := fset.MaxTTLSeconds
		if TTLSeconds <= 0 {
			filter, err = NewSetFilter(s.conn, key, time.Duration(0))
		} else {
			filter, err = NewSetFilter(s.conn, key, time.Duration(TTLSeconds)*time.Second)
		}
	} else {
		fset := setting.GetSetfilterSetting()
		TTLSeconds := fset.MaxTTLSeconds
		if TTLSeconds <= 0 {
			filter, err = NewSetFilter(s.conn, key, time.Duration(0))
		} else {
			filter, err = NewSetFilter(s.conn, key, time.Duration(TTLSeconds)*time.Second)
		}
	}
	if err != nil {
		return err
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

func (s *RedisFilterStore) BlacklistStatus(blacklistid string) (*itemfilterRPC_pb.RedisFilterStatus, error) {
	key := s.BlacklistKey(blacklistid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filter, err := NewSetFilter(s.conn, key, time.Duration(0))
	if err != nil {
		return nil, err
	}
	return filter.Status(ctx)
}
func (s *RedisFilterStore) DeleteBlacklist(blacklistid string) error {
	key := s.BlacklistKey(blacklistid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filter, err := NewSetFilter(s.conn, key, time.Duration(0))
	if err != nil {
		return err
	}
	return filter.Clean(ctx)
}
func (s *RedisFilterStore) UpdateBlacklistSource(blacklistid string, Content ...string) error {
	key := s.BlacklistKey(blacklistid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filter, err := NewSetFilter(s.conn, key, time.Duration(0))
	if err != nil {
		return err
	}
	setting := itemfilterRPC_pb.RedisFilterSetting{
		Setting: &itemfilterRPC_pb.RedisFilterSetting_SetfilterSetting{
			SetfilterSetting: s.Opts.DefaultBlacklistFilterInfo,
		},
	}
	dfset := s.Opts.DefaultBlacklistFilterInfo
	setting.Setting = &itemfilterRPC_pb.RedisFilterSetting_SetfilterSetting{
		SetfilterSetting: dfset,
	}
	err = filter.Clean(ctx)
	if err != nil {
		return err
	}
	err = filter.CreateAndAdd(ctx, &setting, Content...)
	if err != nil {
		return err
	}
	return nil
}
func (s *RedisFilterStore) CheckSourceInBlacklist(blacklistid string, Candidates ...string) (CandidateStatus map[string]bool, err error) {
	key := s.BlacklistKey(blacklistid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filter, err := NewSetFilter(s.conn, key, time.Duration(0))
	if err != nil {
		return nil, err
	}
	return filter.CheckIn(ctx, Candidates...)
}

func (s *RedisFilterStore) FilterSourceInBlacklist(blacklistid string, Need, chunkSize int32, Candidates ...string) (Survivor []string, err error) {
	key := s.BlacklistKey(blacklistid)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filter, err := NewSetFilter(s.conn, key, time.Duration(0))
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
				res := flitout(true, resmap)
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
		return flitout(true, resmap), nil
	}
}
