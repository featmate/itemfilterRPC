package metadatastore

import (
	"strings"
	"time"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"github.com/Golang-Tools/namespace"
	"github.com/Golang-Tools/optparams"
	"github.com/Golang-Tools/redishelper/v2/redisproxy"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

type RedisMetadataStore struct {
	Opts Options
	conn *redisproxy.RedisProxy
}

func NewRedisMetadataStore() *RedisMetadataStore {
	s := new(RedisMetadataStore)
	s.conn = redisproxy.New()
	return s
}

func (s *RedisMetadataStore) Init(opts ...optparams.Option[Options]) error {
	optparams.GetOption(&s.Opts, opts...)
	initparams := []optparams.Option[redisproxy.Options]{}
	address := strings.Split(s.Opts.URL, ",")
	if len(address) > 1 {
		copt := redis.ClusterOptions{
			Addrs: address,
		}
		switch s.Opts.RedisRouteMod {
		case "bylatency":
			{
				copt.RouteByLatency = true
			}
		case "randomly":
			{
				copt.RouteRandomly = true
			}
		}
		initparams = append(initparams, redisproxy.WithClusterOptions(&copt))
	} else {
		initparams = append(initparams, redisproxy.WithURL(s.Opts.URL))
	}
	if s.Opts.ConnTimeout > 0 {
		initparams = append(initparams, redisproxy.WithQueryTimeoutMS(s.Opts.ConnTimeout))
	}
	err := s.conn.Init(initparams...)
	if err != nil {
		return err
	}
	return nil
}

// 黑名单
func (s *RedisMetadataStore) BlacklistRangeKey() string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase}
	key := ns.FullName(s.Opts.BlacklistConfig_RangeInfoKey, namespace.WithRedisStyle())
	return key
}

func (s *RedisMetadataStore) BlacklistInfoNamespace() string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.BlacklistConfig_InfoKeySubnamespace}
	key := ns.ToString(namespace.WithRedisStyle())
	return key
}

func (s *RedisMetadataStore) BlacklistInfoKey(blacklistid string) string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.BlacklistConfig_InfoKeySubnamespace}
	key := ns.FullName(blacklistid, namespace.WithRedisStyle())
	return key
}

func (s *RedisMetadataStore) GetBlacklistList() ([]string, error) {
	rangekey := s.BlacklistRangeKey()
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	return s.conn.SMembers(ctx, rangekey).Result()
}

func (s *RedisMetadataStore) NewBlacklist(info *itemfilterRPC_pb.NewFilterQuery) error {
	rangekey := s.BlacklistRangeKey()
	key := s.BlacklistInfoKey(info.Meta.ID)
	value, err := proto.Marshal(info)
	if err != nil {
		return err
	}
	ttl := time.Duration(0)
	settingMaxTTLSeconds := int64(0)

	if info.Setting.GetSetfilterSetting() != nil {
		settingMaxTTLSeconds = info.Setting.GetSetfilterSetting().MaxTTLSeconds
	}
	if info.Setting.GetBloomfilterSetting() != nil {
		settingMaxTTLSeconds = info.Setting.GetBloomfilterSetting().MaxTTLSeconds
	}
	if info.Setting.GetCuckoofilterSetting() != nil {
		settingMaxTTLSeconds = info.Setting.GetCuckoofilterSetting().MaxTTLSeconds
	}

	if settingMaxTTLSeconds > 0 {
		ttl = time.Duration(settingMaxTTLSeconds) * time.Second
	} else {
		if s.Opts.BlacklistConfig_DefaultTTLDays > 0 {
			hours := s.Opts.BlacklistConfig_DefaultTTLDays * 24
			ttl = time.Duration(int64(hours)) * time.Hour
		}
	}
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	pipe := s.conn.TxPipeline()
	pipe.Set(ctx, key, value, ttl)
	pipe.SAdd(ctx, rangekey, info.Meta.ID)
	_, err = pipe.Exec(ctx)
	return err
}
func (s *RedisMetadataStore) GetBlacklistInfo(id string) (*itemfilterRPC_pb.BlacklistInfo, error) {
	key := s.BlacklistInfoKey(id)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	infostr, err := s.conn.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	}
	info := itemfilterRPC_pb.NewFilterQuery{}
	err = proto.Unmarshal([]byte(infostr), &info)
	if err != nil {
		return nil, err
	}

	result := itemfilterRPC_pb.BlacklistInfo{}
	result.Meta = info.Meta
	result.Setting = info.Setting
	if info.Setting.GetBloomfilterSetting() != nil {
		result.FilterType = itemfilterRPC_pb.RedisFilterType_BLOOM
	} else {
		if info.Setting.GetCuckoofilterSetting() != nil {
			result.FilterType = itemfilterRPC_pb.RedisFilterType_CUCKOO
		} else {
			result.FilterType = itemfilterRPC_pb.RedisFilterType_SET
		}
	}
	return &result, nil
}
func (s *RedisMetadataStore) DeleteBlacklist(id string) error {
	rangekey := s.BlacklistRangeKey()
	key := s.BlacklistInfoKey(id)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	pipe := s.conn.TxPipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, rangekey, id)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// range
func (s *RedisMetadataStore) RangeRangeKey() string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase}
	key := ns.FullName(s.Opts.RangeConfig_RangeInfoKey, namespace.WithRedisStyle())
	return key
}

func (s *RedisMetadataStore) RangeInfoNamespace() string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.RangeConfig_InfoKeySubnamespace}
	key := ns.ToString(namespace.WithRedisStyle())
	return key
}

func (s *RedisMetadataStore) RangeInfoKey(rangeid string) string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.RangeConfig_InfoKeySubnamespace}
	key := ns.FullName(rangeid, namespace.WithRedisStyle())
	return key
}

func (s *RedisMetadataStore) GetRangeList() ([]string, error) {
	rangekey := s.RangeRangeKey()
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	return s.conn.SMembers(ctx, rangekey).Result()
}
func (s *RedisMetadataStore) NewRange(info *itemfilterRPC_pb.NewFilterQuery) error {
	rangekey := s.RangeRangeKey()
	key := s.RangeInfoKey(info.Meta.ID)
	value, err := proto.Marshal(info)
	if err != nil {
		return err
	}
	ttl := time.Duration(0)
	settingMaxTTLSeconds := int64(0)

	if info.Setting.GetSetfilterSetting() != nil {
		settingMaxTTLSeconds = info.Setting.GetSetfilterSetting().MaxTTLSeconds
	}
	if info.Setting.GetBloomfilterSetting() != nil {
		settingMaxTTLSeconds = info.Setting.GetBloomfilterSetting().MaxTTLSeconds
	}
	if info.Setting.GetCuckoofilterSetting() != nil {
		settingMaxTTLSeconds = info.Setting.GetCuckoofilterSetting().MaxTTLSeconds
	}
	if settingMaxTTLSeconds > 0 {
		ttl = time.Duration(settingMaxTTLSeconds) * time.Second
	} else {
		if s.Opts.RangeConfig_DefaultTTLDays > 0 {
			hours := s.Opts.RangeConfig_DefaultTTLDays * 24
			ttl = time.Duration(int64(hours)) * time.Hour
		}
	}
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	pipe := s.conn.TxPipeline()
	pipe.Set(ctx, key, value, ttl)
	pipe.SAdd(ctx, rangekey, info.Meta.ID)
	_, err = pipe.Exec(ctx)
	return err
}
func (s *RedisMetadataStore) GetRangeInfo(id string) (*itemfilterRPC_pb.RangeInfo, error) {
	key := s.RangeInfoKey(id)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	infostr, err := s.conn.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	info := itemfilterRPC_pb.NewFilterQuery{}
	err = proto.Unmarshal([]byte(infostr), &info)
	if err != nil {
		return nil, err
	}
	result := itemfilterRPC_pb.RangeInfo{}
	switch s.Opts.RangeConfig_FilterType {
	case "bloom":
		{
			result.FilterType = itemfilterRPC_pb.RedisFilterType_BLOOM
		}
	case "cuckoo":
		{
			result.FilterType = itemfilterRPC_pb.RedisFilterType_CUCKOO
		}
	default:
		{
			result.FilterType = itemfilterRPC_pb.RedisFilterType_SET
		}
	}

	result.Meta = info.Meta
	result.Setting = info.Setting
	return &result, nil
}
func (s *RedisMetadataStore) DeleteRange(id string) error {
	rangekey := s.RangeRangeKey()
	key := s.RangeInfoKey(id)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	pipe := s.conn.TxPipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, rangekey, id)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
