package metadatastore

import (
	"strings"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"github.com/Golang-Tools/etcdhelper/etcdproxy"
	"github.com/Golang-Tools/namespace"
	"github.com/Golang-Tools/optparams"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

type EtcdMetadataStore struct {
	Opts Options
	conn *etcdproxy.EtcdProxy
}

func NewEtcdMetadataStore() *EtcdMetadataStore {
	s := new(EtcdMetadataStore)
	s.conn = etcdproxy.New()
	return s
}

func (s *EtcdMetadataStore) Init(opts ...optparams.Option[Options]) error {
	optparams.GetOption(&s.Opts, opts...)
	initparams := []optparams.Option[etcdproxy.Options]{}
	if s.Opts.ConnTimeout > 0 {
		initparams = append(initparams, etcdproxy.WithQueryTimeout(s.Opts.ConnTimeout))
	}
	err := s.conn.Init(s.Opts.URL, initparams...)
	if err != nil {
		return err
	}
	return nil
}

// 黑名单
func (s *EtcdMetadataStore) BlacklistRangeKey() string {
	// ns := namespace.NameSpcae{s.Opts.NamespaceBase}
	// key := ns.FullName(s.Opts.BlacklistConfig_RangeInfoKey, namespace.WithEtcdStyle())
	return ""
}

func (s *EtcdMetadataStore) BlacklistInfoNamespace() string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.BlacklistConfig_InfoKeySubnamespace}
	key := ns.ToString(namespace.WithEtcdStyle()) + "/"
	return key
}

func (s *EtcdMetadataStore) BlacklistInfoKey(blacklistid string) string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.BlacklistConfig_InfoKeySubnamespace}
	key := ns.FullName(blacklistid, namespace.WithEtcdStyle())
	return key
}

func (s *EtcdMetadataStore) GetBlacklistList() ([]string, error) {
	prefix := s.BlacklistInfoNamespace()
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	resp, err := s.conn.KV.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		id := strings.ReplaceAll(key, prefix, "")
		res = append(res, id)
	}
	return res, nil
}

func (s *EtcdMetadataStore) NewBlacklist(info *itemfilterRPC_pb.NewFilterQuery) error {
	key := s.BlacklistInfoKey(info.Meta.ID)
	value, err := proto.Marshal(info)
	if err != nil {
		return err
	}
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	putopts := []clientv3.OpOption{}

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
		resp, err := s.conn.Grant(ctx, settingMaxTTLSeconds)
		if err != nil {
			return err
		}
		putopts = append(putopts, clientv3.WithLease(resp.ID))

	} else {
		if s.Opts.BlacklistConfig_DefaultTTLDays > 0 {
			seconds := int64(s.Opts.BlacklistConfig_DefaultTTLDays) * 24 * 3600
			resp, err := s.conn.Grant(ctx, seconds)
			if err != nil {
				return err
			}
			putopts = append(putopts, clientv3.WithLease(resp.ID))
		}
	}
	_, err = s.conn.KV.Put(ctx, key, string(value), putopts...)
	return err
}
func (s *EtcdMetadataStore) GetBlacklistInfo(id string) (*itemfilterRPC_pb.BlacklistInfo, error) {
	key := s.BlacklistInfoKey(id)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	resp, err := s.conn.KV.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, ErrNotFound
	}
	kv := resp.Kvs[0]
	info := itemfilterRPC_pb.NewFilterQuery{}
	err = proto.Unmarshal(kv.Value, &info)
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
func (s *EtcdMetadataStore) DeleteBlacklist(id string) error {
	key := s.BlacklistInfoKey(id)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	_, err := s.conn.KV.Delete(ctx, key)
	return err
}

// range
func (s *EtcdMetadataStore) RangeRangeKey() string {
	// ns := namespace.NameSpcae{s.Opts.NamespaceBase}
	// key := ns.FullName(s.Opts.RangeConfig_RangeInfoKey, namespace.WithEtcdStyle())
	return ""
}

func (s *EtcdMetadataStore) RangeInfoNamespace() string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.RangeConfig_InfoKeySubnamespace}
	key := ns.ToString(namespace.WithEtcdStyle()) + "/"
	return key
}

func (s *EtcdMetadataStore) RangeInfoKey(rangeid string) string {
	ns := namespace.NameSpcae{s.Opts.NamespaceBase, s.Opts.RangeConfig_InfoKeySubnamespace}
	key := ns.FullName(rangeid, namespace.WithEtcdStyle())
	return key
}

func (s *EtcdMetadataStore) GetRangeList() ([]string, error) {
	prefix := s.RangeInfoNamespace()
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	resp, err := s.conn.KV.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		id := strings.ReplaceAll(key, prefix, "")
		res = append(res, id)
	}
	return res, nil
}
func (s *EtcdMetadataStore) NewRange(info *itemfilterRPC_pb.NewFilterQuery) error {
	key := s.RangeInfoKey(info.Meta.ID)
	value, err := proto.Marshal(info)
	if err != nil {
		return err
	}
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	// ttl := time.Duration(0)
	putopts := []clientv3.OpOption{}
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
		resp, err := s.conn.Grant(ctx, settingMaxTTLSeconds)
		if err != nil {
			return err
		}
		putopts = append(putopts, clientv3.WithLease(resp.ID))
	} else {
		if s.Opts.RangeConfig_DefaultTTLDays > 0 {
			seconds := int64(s.Opts.RangeConfig_DefaultTTLDays) * 24 * 3600
			resp, err := s.conn.Grant(ctx, seconds)
			if err != nil {
				return err
			}
			putopts = append(putopts, clientv3.WithLease(resp.ID))
		}
	}
	_, err = s.conn.KV.Put(ctx, key, string(value), putopts...)
	return err
}
func (s *EtcdMetadataStore) GetRangeInfo(id string) (*itemfilterRPC_pb.RangeInfo, error) {
	key := s.RangeInfoKey(id)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	resp, err := s.conn.KV.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, ErrNotFound
	}
	kv := resp.Kvs[0]
	info := itemfilterRPC_pb.NewFilterQuery{}
	err = proto.Unmarshal(kv.Value, &info)
	if err != nil {
		return nil, err
	}
	result := itemfilterRPC_pb.RangeInfo{}
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
func (s *EtcdMetadataStore) DeleteRange(id string) error {
	key := s.RangeInfoKey(id)
	ctx, cancel := s.conn.NewCtx()
	defer cancel()
	_, err := s.conn.KV.Delete(ctx, key)
	return err
}
