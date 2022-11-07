package itemfilterRPC_serv

import (
	"context"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"
)

// GetContextStatus  获取上下文过滤器状态信息
func (s *Server) GetContextStatus(ctx context.Context, in *itemfilterRPC_pb.IDEntitySourceTypeQuery) (*itemfilterRPC_pb.GetContextStatusResponse, error) {
	if in.EntitySourceType == "" {
		return nil, ErrNeedEntityType
	}
	if in.Id == "" {
		return nil, ErrNeedID
	}
	return nil, ErrNotImplemented
}

// NewContexts 批量创建场景
func (s *Server) NewContexts(ctx context.Context, in *itemfilterRPC_pb.NewFilterQuery) (*itemfilterRPC_pb.IDResponse, error) {
	if in.Meta == nil {
		return nil, ErrNeedMetaInfo
	}
	if in.Meta.EntitySourceType == "" {
		return nil, ErrNeedEntityType
	}
	if in.Meta.ID == "" {
		contextid, err := s.idGener.Next()
		if err != nil {
			return nil, err
		}
		in.Meta.ID = contextid
	} else {
		status, err := s.filterstore.ContextStatus(in.Meta.EntitySourceType, in.Meta.ID)
		if err != nil {
			return nil, err
		}
		if status.Exists {
			return nil, ErrIDExists
		}
	}
	err := s.filterstore.NewContext(in.Meta.EntitySourceType, in.Meta.ID, in.Setting, in.EntitySources...)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.IDResponse{
		ID: in.Meta.ID,
	}, nil
}
