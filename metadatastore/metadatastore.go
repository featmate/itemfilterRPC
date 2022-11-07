package metadatastore

import (
	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"github.com/Golang-Tools/optparams"
)

type MetadataStore interface {
	Init(opts ...optparams.Option[Options]) error
	// 黑名单
	BlacklistRangeKey() string
	BlacklistInfoNamespace() string
	BlacklistInfoKey(blacklistid string) string

	GetBlacklistList() ([]string, error)
	NewBlacklist(*itemfilterRPC_pb.NewFilterQuery) error
	GetBlacklistInfo(id string) (*itemfilterRPC_pb.BlacklistInfo, error)
	DeleteBlacklist(id string) error

	// range
	RangeRangeKey() string
	RangeInfoNamespace() string
	RangeInfoKey(rangeid string) string

	GetRangeList() ([]string, error)
	NewRange(*itemfilterRPC_pb.NewFilterQuery) error
	GetRangeInfo(id string) (*itemfilterRPC_pb.RangeInfo, error)
	DeleteRange(id string) error
}
