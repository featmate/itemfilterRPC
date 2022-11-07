package itemfilterRPC_serv

import "errors"

var ErrNotImplemented = errors.New("NotImplemented")

var ErrNotFound = errors.New("NotFound")

var ErrIDExists = errors.New("IDExists")

var ErrNeedID = errors.New("NeedID")

var ErrNeedEntityType = errors.New("NeedEntityType")

var ErrNeedMetaInfo = errors.New("NeedMetaInfo")

var ErrNeedAtLeastOneCandidate = errors.New("NeedAtLeastOneCandidate")

var ErrNeedConditions = errors.New("NeedConditions")

var ErrNeedAtLeastOneEntitySource = errors.New("NeedAtLeastOneEntitySource")
