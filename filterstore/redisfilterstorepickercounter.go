package filterstore

import (
	"time"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"
	"github.com/go-redis/redis/v8"

	"github.com/Golang-Tools/namespace"
)

// pickercounter
func (s *RedisFilterStore) PickerCounterKeyTemplate() string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.NamespaceConfig_PickerCounterSubnamespace, s.Opts.NamespaceConfig_PickerCounterDateSubnamespace, "{{ YYYY-MM-DD }}"}
	key := ns.ToString(namespace.WithRedisStyle())
	return key
}
func (s *RedisFilterStore) PickerCounterKey(Date time.Time) string {
	datastr := Date.Format("2006-01-02")
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.NamespaceConfig_PickerCounterSubnamespace, s.Opts.NamespaceConfig_PickerCounterDateSubnamespace, datastr}
	key := ns.ToString(namespace.WithRedisStyle())
	return key
}
func (s *RedisFilterStore) RangePickerCounterKeyTemplate() string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.NamespaceConfig_PickerCounterSubnamespace, s.Opts.NamespaceConfig_PickerCounterDateSubnamespace, "{{ YYYY-MM-DD }}", s.Opts.RangeConfig_KeySubnamespace, "{{ RID }}"}
	key := ns.ToString(namespace.WithRedisStyle())
	return key
}
func (s *RedisFilterStore) RangePickerCounterKey(Date time.Time, RID string) string {
	datastr := Date.Format("2006-01-02")
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.NamespaceConfig_PickerCounterSubnamespace, s.Opts.NamespaceConfig_PickerCounterDateSubnamespace, datastr, s.Opts.RangeConfig_KeySubnamespace, RID}
	key := ns.ToString(namespace.WithRedisStyle())
	return key
}

func (s *RedisFilterStore) GetPickerCounterNumber(Days, Offset int32, mod itemfilterRPC_pb.UnionMod, Ranges ...string) (map[string]int64, error) {
	if !s.Opts.PickerCounterOn {
		return nil, ErrFeatureNotOn
	}

	if len(Ranges) > 0 {
		if !s.Opts.PickerCounterConfig_RangeCounterOn {
			return nil, ErrFeatureNotOn
		} else {
			if Days+Offset > int32(s.Opts.PickerCounterConfig_RangeTTLDays) {
				return nil, ErrDateOversize
			}
		}
	} else {
		if Days+Offset > int32(s.Opts.PickerCounterConfig_TTLDays) {
			return nil, ErrDateOversize
		}
	}

	dates := []time.Time{}
	starttime := time.Now().In(time.UTC).AddDate(0, 0, -int(Offset))
	dates = append(dates, starttime)
	for i := 0; i < int(Days); i++ {
		dates = append(dates, starttime.AddDate(0, 0, -int(i)))
	}

	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	pipe := s.conn.TxPipeline()
	cmds := map[string]*redis.IntCmd{}
	if len(Ranges) > 0 {
		switch mod {
		case itemfilterRPC_pb.UnionMod_UnionByRange:
			{
				for _, date := range dates {
					keys := []string{}
					for _, rangeid := range Ranges {
						key := s.RangePickerCounterKey(date, rangeid)
						keys = append(keys, key)
					}
					datastr := date.Format("2006-01-02")
					cmds[datastr] = pipe.PFCount(ctx, keys...)
				}
			}
		case itemfilterRPC_pb.UnionMod_UnionByDay:
			{
				for _, rangeid := range Ranges {
					keys := []string{}
					for _, date := range dates {
						key := s.RangePickerCounterKey(date, rangeid)
						keys = append(keys, key)
					}
					cmds[rangeid] = pipe.PFCount(ctx, keys...)
				}
			}
		default:
			{
				keys := []string{}
				for _, date := range dates {
					for _, rangeid := range Ranges {
						key := s.RangePickerCounterKey(date, rangeid)
						keys = append(keys, key)
					}
				}
				cmds["union"] = pipe.PFCount(ctx, keys...)
			}
		}
	} else {
		switch mod {
		case itemfilterRPC_pb.UnionMod_UnionByRange:
			{

				for _, date := range dates {
					key := s.PickerCounterKey(date)
					datastr := date.Format("2006-01-02")
					cmds[datastr] = pipe.PFCount(ctx, key)
				}
			}
		default:
			{
				keys := []string{}
				for _, date := range dates {
					key := s.PickerCounterKey(date)
					keys = append(keys, key)
				}
				cmds["global"] = pipe.PFCount(ctx, keys...)
			}

		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	result := map[string]int64{}
	for key, cmd := range cmds {
		r, err := cmd.Result()
		if err != nil {
			return nil, err
		}
		result[key] = r
	}
	return result, ErrNotImplemented
}
