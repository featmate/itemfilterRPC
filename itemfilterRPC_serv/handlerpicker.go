package itemfilterRPC_serv

import (
	"context"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"
)

// GetPickerStatus  获取上下文过滤器状态信息
func (s *Server) GetPickerStatus(ctx context.Context, in *itemfilterRPC_pb.IDEntitySourceTypeQuery) (*itemfilterRPC_pb.GetPickerStatusResponse, error) {
	if in.EntitySourceType == "" {
		return nil, ErrNeedEntityType
	}
	if in.Id == "" {
		return nil, ErrNeedID
	}
	r, err := s.filterstore.PickerStatus(in.EntitySourceType, in.Id)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.GetPickerStatusResponse{
		FilterStatus: r,
	}, nil
}

// NewPickers 批量创建Picker
func (s *Server) NewPickers(ctx context.Context, in *itemfilterRPC_pb.NewPickerFilterQuery) (*itemfilterRPC_pb.IDResponse, error) {

	if in.Meta == nil {
		return nil, ErrNeedMetaInfo
	}
	if in.Meta.EntitySourceType == "" {
		return nil, ErrNeedEntityType
	}
	if in.Meta.ID == "" {
		pickerid, err := s.idGener.Next()
		if err != nil {
			return nil, err
		}
		in.Meta.ID = pickerid
	}
	err := s.filterstore.NewPicker(in.Meta.EntitySourceType, in.Meta.ID, in.Settings, in.EntitySources...)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.IDResponse{ID: in.Meta.ID}, nil
}
