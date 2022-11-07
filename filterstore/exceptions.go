package filterstore

import "errors"

var ErrNotImplemented = errors.New("NotImplemented")

var ErrFilterSettingNotMatchFilterType = errors.New("FilterSettingNotMatchFilterType")
var ErrKeyNotFilter = errors.New("KeyNotFilter")

var ErrNeedAtLeastOneEnvInfo = errors.New("NeedAtLeastOneEnvInfo")

var ErrFeatureNotOn = errors.New("FeatureNotOn")

var ErrKeyNotExist = errors.New("KeyNotExist")

var ErrDateOversize = errors.New("DateOversize")

var ErrRangeNotExist = errors.New("RangeNotExist")
