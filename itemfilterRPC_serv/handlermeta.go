package itemfilterRPC_serv

import (
	"context"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

// GetMeta 获取本服务的元信息
func (s *Server) GetMeta(ctx context.Context, in *emptypb.Empty) (*itemfilterRPC_pb.MetaResponse, error) {
	var MetaDataStoreType itemfilterRPC_pb.MetaDataStoreType
	switch s.MetaDataStore_Type {
	case "redis":
		{
			MetaDataStoreType = itemfilterRPC_pb.MetaDataStoreType_REDIS
		}
	case "etcd":
		{
			MetaDataStoreType = itemfilterRPC_pb.MetaDataStoreType_ETCD
		}
	default:
		{
			MetaDataStoreType = itemfilterRPC_pb.MetaDataStoreType_SELF
		}

	}

	m := &itemfilterRPC_pb.MetaResponse{
		Content: &itemfilterRPC_pb.MetaInfo{
			RedisURL: s.FilterRedis_URL,
			MetaDataStoreConfig: &itemfilterRPC_pb.MetaMetaDataStoreConfig{
				StoreType: MetaDataStoreType,
				StoreURL:  s.MetaDataStore_URL,
			},
			BlacklistConfig: &itemfilterRPC_pb.MetaBlacklistConfig{
				RangeInfoKey:       s.metadatastore.BlacklistRangeKey(),
				InfoKeyNamespace:   s.metadatastore.BlacklistInfoNamespace(),
				InfoDefaultTTLDays: int32(s.BlacklistConfig_DefaultTTLDays),
			},
			RangeConfig: &itemfilterRPC_pb.MetaRangeConfig{
				RangeInfoKey:       s.metadatastore.RangeRangeKey(),
				InfoKeyNamespace:   s.metadatastore.RangeInfoNamespace(),
				InfoDefaultTTLDays: int32(s.RangeConfig_DefaultTTLDays),
			},
			ContextConfig: &itemfilterRPC_pb.MetaContextConfig{
				KeyTemplate:    s.filterstore.ContextKeyTemplate(),
				DefaultSetting: s.defaultContextFilterInfo,
			},
			PickerConfig: &itemfilterRPC_pb.MetaPickerConfig{
				KeyTemplate: s.filterstore.PickerKeyTemplate(),
				FilterInfos: s.defaultPickerFilterInfos,
			},
			PickerCounterOn: s.PickerCounterOn,
			PickerCounterConfig: &itemfilterRPC_pb.MetaPickerCounterConfig{
				TTLDays:                       int32(s.PickerCounterConfig_TTLDays),
				PickerCounterKeyTemplate:      s.filterstore.PickerCounterKeyTemplate(),
				RangeCounterOn:                s.PickerCounterConfig_RangeCounterOn,
				RangeTTLDays:                  int32(s.PickerCounterConfig_RangeTTLDays),
				RangePickerCounterKeyTemplate: s.filterstore.RangePickerCounterKeyTemplate(),
			},
		},
	}
	return m, nil
}
