package filterstore

import (
	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"github.com/Golang-Tools/optparams"
)

// Option 设置key行为的选项
type Options struct {
	RedisURL                                      string
	RedisRouteMod                                 string
	ConnTimeout                                   int
	NamespaceBase                                 string
	BlacklistConfig_KeySubnamespace               string
	BlacklistConfig_DefaultTTLDays                int
	RangeConfig_KeySubnamespace                   string
	RangeConfig_FilterType                        string
	NamespaceConfig_EntitySourceTypeSubnamespace  string
	ContextConfig_KeySubnamespace                 string
	NamespaceConfig_PickerIntervalSubnamespace    string
	NamespaceConfig_PickerSubnamespace            string
	NamespaceConfig_PickerCounterSubnamespace     string
	NamespaceConfig_PickerCounterDateSubnamespace string

	PickerCounterOn                    bool
	PickerCounterConfig_TTLDays        int
	PickerCounterConfig_RangeCounterOn bool
	PickerCounterConfig_RangeTTLDays   int

	DefaultBlacklistFilterInfo *itemfilterRPC_pb.SetFilterSetting
	DefaultRangeFilterInfo     *itemfilterRPC_pb.RedisFilterSetting
	DefaultContextFilterInfo   *itemfilterRPC_pb.RedisFilterSetting
	DefaultPickerFilterInfos   map[string]*itemfilterRPC_pb.RedisFilterSetting
}

var DefaultOpts = Options{
	RedisURL:                        "redis://localhost:6379/0",
	RedisRouteMod:                   "",
	ConnTimeout:                     50,
	NamespaceBase:                   "itemfilter",
	BlacklistConfig_KeySubnamespace: "BID",
	RangeConfig_KeySubnamespace:     "RID",
	NamespaceConfig_EntitySourceTypeSubnamespace:  "EST",
	ContextConfig_KeySubnamespace:                 "CID",
	NamespaceConfig_PickerIntervalSubnamespace:    "PI",
	NamespaceConfig_PickerSubnamespace:            "PID",
	NamespaceConfig_PickerCounterSubnamespace:     "PC",
	NamespaceConfig_PickerCounterDateSubnamespace: "PCD",
	PickerCounterOn:                    false,
	PickerCounterConfig_TTLDays:        7,
	PickerCounterConfig_RangeCounterOn: false,
	PickerCounterConfig_RangeTTLDays:   7,
}

// WithOptions 使用特定单机redis连接设置
func WithOptions(opts *Options) optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.RedisURL = opts.RedisURL
		o.RedisRouteMod = opts.RedisRouteMod
		o.NamespaceBase = opts.NamespaceBase
		o.BlacklistConfig_KeySubnamespace = opts.BlacklistConfig_KeySubnamespace
		o.BlacklistConfig_DefaultTTLDays = opts.BlacklistConfig_DefaultTTLDays
		o.RangeConfig_KeySubnamespace = opts.RangeConfig_KeySubnamespace
		o.RangeConfig_FilterType = opts.RangeConfig_FilterType
		o.NamespaceConfig_EntitySourceTypeSubnamespace = opts.NamespaceConfig_EntitySourceTypeSubnamespace
		o.ContextConfig_KeySubnamespace = opts.ContextConfig_KeySubnamespace
		o.NamespaceConfig_PickerIntervalSubnamespace = opts.NamespaceConfig_PickerIntervalSubnamespace
		o.NamespaceConfig_PickerSubnamespace = opts.NamespaceConfig_PickerSubnamespace
		o.NamespaceConfig_PickerCounterSubnamespace = opts.NamespaceConfig_PickerCounterSubnamespace
		o.NamespaceConfig_PickerCounterDateSubnamespace = opts.NamespaceConfig_PickerCounterDateSubnamespace
		o.ConnTimeout = opts.ConnTimeout
		o.PickerCounterOn = opts.PickerCounterOn
		o.PickerCounterConfig_TTLDays = opts.PickerCounterConfig_TTLDays
		o.PickerCounterConfig_RangeCounterOn = opts.PickerCounterConfig_RangeCounterOn
		o.PickerCounterConfig_RangeTTLDays = opts.PickerCounterConfig_RangeTTLDays
		o.DefaultBlacklistFilterInfo = opts.DefaultBlacklistFilterInfo
		o.DefaultRangeFilterInfo = opts.DefaultRangeFilterInfo
		o.DefaultContextFilterInfo = opts.DefaultContextFilterInfo
		o.DefaultPickerFilterInfos = opts.DefaultPickerFilterInfos
	})
}
