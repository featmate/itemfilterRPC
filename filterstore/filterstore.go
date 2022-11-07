package filterstore

import (
	"time"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"github.com/Golang-Tools/optparams"
)

type FilterStore interface {
	Init(opts ...optparams.Option[Options]) error

	//黑名单
	BlacklistKey(blacklistid string) string

	NewBlacklist(blacklistid string, setting *itemfilterRPC_pb.RedisFilterSetting, Content ...string) error
	BlacklistStatus(blacklistid string) (*itemfilterRPC_pb.RedisFilterStatus, error)
	DeleteBlacklist(id string) error
	UpdateBlacklistSource(id string, Content ...string) error
	CheckSourceInBlacklist(id string, Candidates ...string) (CandidateStatus map[string]bool, err error)
	FilterSourceInBlacklist(id string, Need, ChunkSize int32, Candidates ...string) (Survivor []string, err error)

	//range
	RangeKey(rangeid string) string

	NewRange(rangeid string, setting *itemfilterRPC_pb.RedisFilterSetting, Content ...string) error
	RangeStatus(rangeid string) (*itemfilterRPC_pb.RedisFilterStatus, error)
	DeleteRange(id string) error
	AddRangeSource(id string, Content ...string) error
	CheckSourceInRange(id string, Candidates ...string) (CandidateStatus map[string]bool, err error)
	FilterSourceNotInRange(id string, Need, ChunkSize int32, Candidates ...string) (Survivor []string, err error)

	//context
	ContextKeyTemplate() string
	ContextKey(EntitySourceType, ContextID string) string
	NewContext(EntitySourceType, ContextID string, setting *itemfilterRPC_pb.RedisFilterSetting, Content ...string) error
	ContextStatus(EntitySourceType, ContextID string) (*itemfilterRPC_pb.RedisFilterStatus, error)

	//picker
	PickerKeyTemplate() string
	PickerKey(EntitySourceType, PickerInterval, PickerID string) string
	NewPicker(EntitySourceType, PickerID string, settings map[string]*itemfilterRPC_pb.RedisFilterSetting, Content ...string) error
	PickerStatus(EntitySourceType, PickerID string) (map[string]*itemfilterRPC_pb.RedisFilterStatus, error)
	//opra
	SetEntitySourceUsed(*itemfilterRPC_pb.SetSourceUsedQuery) error
	CheckEntitySource(Condition *itemfilterRPC_pb.Condition, Candidates ...string) (CandidateStatus map[string]bool, err error)
	FilterEntitySource(Condition *itemfilterRPC_pb.Condition, Need, ChunkSize int32, Candidates ...string) (Survivor []string, err error)

	//pickercounter
	PickerCounterKeyTemplate() string
	PickerCounterKey(Date time.Time) string
	RangePickerCounterKeyTemplate() string
	RangePickerCounterKey(Date time.Time, RID string) string

	GetPickerCounterNumber(Days, Offset int32, mod itemfilterRPC_pb.UnionMod, Ranges ...string) (map[string]int64, error)
}
