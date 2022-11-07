package metadatastore

import (
	"github.com/Golang-Tools/optparams"
)

// Option 设置key行为的选项
type Options struct {
	URL                                 string
	RedisRouteMod                       string
	ConnTimeout                         int
	NamespaceBase                       string
	BlacklistConfig_RangeInfoKey        string
	BlacklistConfig_InfoKeySubnamespace string
	RangeConfig_RangeInfoKey            string
	RangeConfig_InfoKeySubnamespace     string
	RangeConfig_FilterType              string
	BlacklistConfig_DefaultTTLDays      int
	RangeConfig_DefaultTTLDays          int
}

var DefaultOpts = Options{
	ConnTimeout:                         50,
	NamespaceBase:                       "itemfilter",
	BlacklistConfig_RangeInfoKey:        "blacklists",
	BlacklistConfig_InfoKeySubnamespace: "infoblacklist",
	RangeConfig_RangeInfoKey:            "range",
	RangeConfig_InfoKeySubnamespace:     "inforange",
	RangeConfig_FilterType:              "set",
}

// WithOptions 使用特定单机redis连接设置
func WithOptions(opts *Options) optparams.Option[Options] {
	return optparams.NewFuncOption(func(o *Options) {
		o.URL = opts.URL
		o.RedisRouteMod = opts.RedisRouteMod
		o.ConnTimeout = opts.ConnTimeout
		o.NamespaceBase = opts.NamespaceBase
		o.BlacklistConfig_RangeInfoKey = opts.BlacklistConfig_RangeInfoKey
		o.BlacklistConfig_InfoKeySubnamespace = opts.BlacklistConfig_InfoKeySubnamespace
		o.RangeConfig_RangeInfoKey = opts.RangeConfig_RangeInfoKey
		o.RangeConfig_InfoKeySubnamespace = opts.RangeConfig_InfoKeySubnamespace
		o.RangeConfig_FilterType = opts.RangeConfig_FilterType
		o.BlacklistConfig_DefaultTTLDays = opts.BlacklistConfig_DefaultTTLDays
		o.RangeConfig_DefaultTTLDays = opts.RangeConfig_DefaultTTLDays
	})
}
