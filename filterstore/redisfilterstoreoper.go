package filterstore

import (
	"context"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"
)

// opr

func (s *RedisFilterStore) SetEntitySourceUsed(in *itemfilterRPC_pb.SetSourceUsedQuery) error {
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	filterPromises := []*FilterPromise{}
	pipe := s.conn.TxPipeline()
	if in.EnvInfo.PickerId == "" {
		if len(in.EnvInfo.ContextIds) == 0 {
			return ErrNeedAtLeastOneEnvInfo
		}
		for contextid, ttlseconds := range in.EnvInfo.ContextIds {
			key := s.ContextKey(in.EnvInfo.EntitySourceType, contextid)
			filterPromises = append(filterPromises, findFilterFromKeyPipe(pipe, ctx, key, time.Duration(ttlseconds)*time.Second, "", s.Opts.DefaultContextFilterInfo))
		}

	} else {
		for TimeRange, setting := range s.Opts.DefaultPickerFilterInfos {
			key := s.PickerKey(in.EnvInfo.EntitySourceType, TimeRange, in.EnvInfo.PickerId)
			filterPromises = append(filterPromises, findFilterFromKeyPipe(pipe, ctx, key, time.Duration(0), TimeRange, setting))
		}
		if len(in.EnvInfo.ContextIds) > 0 {
			for contextid, ttlseconds := range in.EnvInfo.ContextIds {
				key := s.ContextKey(in.EnvInfo.EntitySourceType, contextid)
				filterPromises = append(filterPromises, findFilterFromKeyPipe(pipe, ctx, key, time.Duration(ttlseconds)*time.Second, "", s.Opts.DefaultContextFilterInfo))
			}
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	pipe = s.conn.TxPipeline()
	for _, filterpromise := range filterPromises {
		filter, err := filterpromise.Result(s.conn)
		if err != nil {
			return err
		}
		err = filter.PipeCreateAndAdd(pipe, ctx, filterpromise.Defaultsetting, filterpromise.Expat, in.EntitySources...)
		if err != nil {
			return err
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *RedisFilterStore) checkEntitySourecChunk(ctx context.Context, Condition *itemfilterRPC_pb.Condition, candidates ...string) (map[string]bool, []map[string]bool, error) {
	filterPromises := []*FilterPromise{}
	var rangeFilterPromise *FilterPromise
	var mustIn map[string]bool
	mustNotIn := []map[string]bool{}
	pipe := s.conn.TxPipeline()
	now := time.Now().In(time.UTC)
	if Condition.PickerId != "" {
		if s.Opts.PickerCounterOn {
			key := s.PickerCounterKey(now)
			pipe.PFAdd(ctx, key, Condition.PickerId)
			pipe.Expire(ctx, key, time.Duration(s.Opts.PickerCounterConfig_TTLDays)*24*time.Hour)
		}
		for TimeRange := range s.Opts.DefaultPickerFilterInfos {
			key := s.PickerKey(Condition.EntitySourceType, TimeRange, Condition.PickerId)
			filterPromises = append(filterPromises, findFilterFromKeyPipe(pipe, ctx, key, time.Duration(0), "", nil))
		}
	}
	if len(Condition.ContextIds) > 0 {
		for _, contextid := range Condition.ContextIds {
			key := s.ContextKey(Condition.EntitySourceType, contextid)
			filterPromises = append(filterPromises, findFilterFromKeyPipe(pipe, ctx, key, time.Duration(0), "", nil))
		}
	}
	if len(Condition.BlacklistIds) > 0 {
		for _, blacklistid := range Condition.BlacklistIds {
			key := s.BlacklistKey(blacklistid)
			filterPromises = append(filterPromises, findFilterFromKeyPipe(pipe, ctx, key, time.Duration(0), "", nil))
		}
	}

	if Condition.RangeId != "" {
		key := s.RangeKey(Condition.RangeId)
		rangeFilterPromise = findFilterFromKeyPipe(pipe, ctx, key, time.Duration(0), "", nil)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, nil, err
	}
	// 处理range
	if rangeFilterPromise != nil {
		pipe = s.conn.TxPipeline()
		f, err := rangeFilterPromise.Result(s.conn)
		if err != nil {
			if err == ErrKeyNotExist {
				return nil, nil, ErrRangeNotExist
			} else {
				return nil, nil, err
			}
		}
		rangresultpromise, err := f.PipeCheckIn(pipe, ctx, candidates...)
		if err != nil {
			return nil, nil, err
		}
		if s.Opts.PickerCounterOn && Condition.PickerId != "" && s.Opts.PickerCounterConfig_RangeCounterOn {
			key := s.RangePickerCounterKey(now, Condition.RangeId)
			pipe.PFAdd(ctx, key, Condition.PickerId)
			pipe.Expire(ctx, key, time.Duration(s.Opts.PickerCounterConfig_RangeTTLDays)*24*time.Hour)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			return nil, nil, err
		}
		mustIn, err = rangresultpromise.Result()
		if err != nil {
			return nil, nil, err
		}
	}
	// 处理其他
	if len(filterPromises) > 0 {
		pipe = s.conn.TxPipeline()
		rpromise := []BoolMapPromise{}
		for _, filterPromise := range filterPromises {
			f, err := filterPromise.Result(s.conn)
			if err != nil {
				if err == ErrKeyNotExist {
					continue
				} else {
					return nil, nil, err
				}
			}
			r, err := f.PipeCheckIn(pipe, ctx, candidates...)
			if err != nil {
				return nil, nil, err
			}
			rpromise = append(rpromise, r)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			return nil, nil, err
		}
		for _, rp := range rpromise {
			r, err := rp.Result()
			if err != nil {
				return nil, nil, err
			}
			mustNotIn = append(mustNotIn, r)
		}
	}
	return mustIn, mustNotIn, nil
}

// CheckEntitySource 如果range不存在则报错,其他条件如果不存在则跳过校验
func (s *RedisFilterStore) CheckEntitySource(Condition *itemfilterRPC_pb.Condition, Candidates ...string) (CandidateStatus map[string]bool, err error) {
	mySet := mapset.NewSet(Candidates...)
	candidates := mySet.ToSlice()
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	mustIn, mustNotIn, err := s.checkEntitySourecChunk(ctx, Condition, candidates...)
	if err != nil {
		return nil, err
	}
	result := map[string]bool{}
	for _, cand := range candidates {
		result[cand] = true
		if len(mustIn) > 0 {
			in, ok := mustIn[cand]
			if !ok || !in {
				result[cand] = false
				continue
			}
		}
		if len(mustNotIn) > 0 {
			for _, notinmap := range mustNotIn {
				in, ok := notinmap[cand]
				if ok && in {
					result[cand] = false
					break
				}
			}
		}
	}
	return result, nil
}
func (s *RedisFilterStore) FilterEntitySource(Condition *itemfilterRPC_pb.Condition, Need, chunkSize int32, Candidates ...string) (Survivor []string, err error) {
	mySet := mapset.NewSet(Candidates...)
	candidates := mySet.ToSlice()
	ctx, cancel := s.conn.NewCtx()
	defer cancel()

	if Need > 0 {
		if chunkSize < 1 {
			mustIn, mustNotIn, err := s.checkEntitySourecChunk(ctx, Condition, candidates...)
			if err != nil {
				return nil, err
			}
			result := []string{}
			for _, cand := range candidates {
				if len(result) > int(Need) {
					return result[:Need], nil
				} else {
					if len(mustIn) > 0 {
						in, ok := mustIn[cand]
						if !ok || !in {
							continue
						}
					}
					if len(mustNotIn) > 0 {
						app := true
						for _, notinmap := range mustNotIn {
							in, ok := notinmap[cand]
							if ok && in {
								app = false
								break
							}
						}
						if app {
							result = append(result, cand)
						}
					} else {
						result = append(result, cand)
					}
				}
			}
			return result, nil
		} else {
			chunks := chunkBy(candidates, int(chunkSize))
			result := []string{}
			for _, chunk := range chunks {
				mustIn, mustNotIn, err := s.checkEntitySourecChunk(ctx, Condition, chunk...)
				if err != nil {
					return nil, err
				}
				for _, cand := range chunk {
					if len(result) >= int(Need) {
						return result, nil
					} else {
						if len(mustIn) > 0 {
							in, ok := mustIn[cand]
							if !ok || !in {
								continue
							}
						}
						if len(mustNotIn) > 0 {
							app := true
							for _, notinmap := range mustNotIn {
								in, ok := notinmap[cand]
								if ok && in {
									app = false
									break
								}
							}
							if app {
								result = append(result, cand)
							}
						} else {
							result = append(result, cand)
						}
					}
				}
			}
			return result, nil
		}

	} else {
		mustIn, mustNotIn, err := s.checkEntitySourecChunk(ctx, Condition, candidates...)
		if err != nil {
			return nil, err
		}
		result := []string{}
		for _, cand := range candidates {
			if len(mustIn) > 0 {
				in, ok := mustIn[cand]
				if !ok || !in {
					continue
				}
			}
			if len(mustNotIn) > 0 {
				app := true
				for _, notinmap := range mustNotIn {
					in, ok := notinmap[cand]
					if ok && in {
						app = false
						break
					}
				}
				if app {
					result = append(result, cand)
				}
			} else {
				result = append(result, cand)
			}
		}
		return result, nil
	}
}
