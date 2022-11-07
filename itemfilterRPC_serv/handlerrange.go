package itemfilterRPC_serv

import (
	"context"
	"io"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

// GetRangeList 获取范围列表
func (s *Server) GetRangeList(ctx context.Context, in *emptypb.Empty) (*itemfilterRPC_pb.ListResponse, error) {
	r, err := s.metadatastore.GetRangeList()
	if err != nil {
		return nil, err
	}
	result := itemfilterRPC_pb.ListResponse{
		Content: r,
	}
	return &result, nil
}

// NewRange 创建一个新的range,注意创建并不会向其中添加元素,添加元素请调用AddRangeSource接口
func (s *Server) NewRange(ctx context.Context, in *itemfilterRPC_pb.NewFilterQuery) (*itemfilterRPC_pb.IDResponse, error) {
	if in.Meta == nil {
		return nil, ErrNeedMetaInfo
	}
	if in.Meta.ID == "" {
		_rangeid, err := s.idGener.Next()
		if err != nil {
			return nil, err
		}
		in.Meta.ID = _rangeid
	} else {
		status, err := s.filterstore.RangeStatus(in.Meta.ID)
		if err != nil {
			return nil, err
		}
		if status.Exists {
			return nil, ErrIDExists
		}
	}
	in.Meta.Key = s.filterstore.RangeKey(in.Meta.ID)
	err := s.filterstore.NewRange(in.Meta.ID, in.Setting, in.EntitySources...)
	if err != nil {
		return nil, err
	}
	err = s.metadatastore.NewRange(in)
	if err != nil {
		//回退
		s.filterstore.DeleteRange(in.Meta.ID)
		return nil, err
	}
	return &itemfilterRPC_pb.IDResponse{
		ID: in.Meta.ID,
	}, nil
}

// GetRangeInfo 获取范围信息
func (s *Server) GetRangeInfo(ctx context.Context, in *itemfilterRPC_pb.IDQuery) (*itemfilterRPC_pb.RangeInfoResponse, error) {
	if in.Id == "" {
		return nil, ErrNeedID
	}
	rangeid := in.Id
	status, err := s.filterstore.RangeStatus(rangeid)
	if err != nil {
		return nil, err
	}
	if !status.Exists {
		return nil, ErrNotFound
	}
	info, err := s.metadatastore.GetRangeInfo(rangeid)
	if err != nil {
		return nil, err
	}
	info.Status = status
	return &itemfilterRPC_pb.RangeInfoResponse{Info: info}, nil
}

// deleteRange 指定id后尝试删除这个id对应的信息
func (s *Server) deleteRange(rangeid string) error {
	err := s.filterstore.DeleteRange(rangeid)
	if err != nil {
		return err
	}
	err = s.metadatastore.DeleteRange(rangeid)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRange 删除范围
func (s *Server) DeleteRange(ctx context.Context, in *itemfilterRPC_pb.IDQuery) (*itemfilterRPC_pb.IDResponse, error) {
	if in.Id == "" {
		return nil, ErrNeedID
	}
	rangeid := in.Id
	err := s.deleteRange(rangeid)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.IDResponse{
		ID: rangeid,
	}, nil
}

// AddRangeSource 添加范围内元素
func (s *Server) AddRangeSource(ctx context.Context, in *itemfilterRPC_pb.AddRangeSourceQuery) (*itemfilterRPC_pb.IDResponse, error) {
	if in.Id == "" {
		return nil, ErrNeedID
	}
	rangeid := in.Id
	status, err := s.filterstore.RangeStatus(rangeid)
	if err != nil {
		return nil, err
	}
	if !status.Exists {
		return nil, ErrNotFound
	}
	err = s.filterstore.AddRangeSource(rangeid, in.Content...)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.IDResponse{ID: rangeid}, nil
}

// BatchAddRangeSource 添加范围内元素
func (s *Server) BatchAddRangeSource(stream itemfilterRPC_pb.ITEMFILTERRPC_BatchAddRangeSourceServer) error {
	succeed := []string{}
	rangeid := ""
	for {
		in, err := stream.Recv()

		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&itemfilterRPC_pb.ListResponse{
					Content: succeed,
				})
			}
			return err
		}
		if rangeid == "" {
			if in.Id == "" {
				return ErrNeedID
			} else {
				rangeid = in.Id
				status, err := s.filterstore.RangeStatus(rangeid)
				if err != nil {
					return err
				}
				if !status.Exists {
					return ErrNotFound
				}
			}
		}
		err = s.filterstore.AddRangeSource(rangeid, in.Content...)
		if err != nil {
			return err
		}
		succeed = append(succeed, rangeid)
	}
}

// CheckSourceInRange 检查物品在不在range中
func (s *Server) CheckSourceInRange(ctx context.Context, in *itemfilterRPC_pb.CheckSourceQuery) (*itemfilterRPC_pb.CandidateStatusResponse, error) {
	if in.Id == "" {
		return nil, ErrNeedID
	}
	if len(in.Candidates) == 0 {
		return nil, ErrNeedAtLeastOneCandidate
	}

	rangeid := in.Id
	status, err := s.filterstore.RangeStatus(rangeid)
	if err != nil {
		return nil, err
	}
	if !status.Exists {
		return nil, ErrNotFound
	}
	info, err := s.metadatastore.GetRangeInfo(rangeid)
	if err != nil {
		return nil, err
	}
	res, err := s.filterstore.CheckSourceInRange(rangeid, in.Candidates...)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.CandidateStatusResponse{
		EntitySourceType: info.Meta.EntitySourceType,
		CandidateStatus:  res,
	}, nil
}

// FilterSourceNotInRange 过滤掉不在范围内的物品,保留在的部分
func (s *Server) FilterSourceNotInRange(ctx context.Context, in *itemfilterRPC_pb.FilterSourceQuery) (*itemfilterRPC_pb.SurvivorResponse, error) {
	if in.Id == "" {
		return nil, ErrNeedID
	}
	if len(in.Candidates) == 0 {
		return nil, ErrNeedAtLeastOneCandidate
	}

	rangeid := in.Id
	status, err := s.filterstore.RangeStatus(rangeid)
	if err != nil {
		return nil, err
	}
	if !status.Exists {
		return nil, ErrNotFound
	}
	info, err := s.metadatastore.GetRangeInfo(rangeid)
	if err != nil {
		return nil, err
	}
	res, err := s.filterstore.FilterSourceNotInRange(rangeid, in.Need, in.ChunkSize, in.Candidates...)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.SurvivorResponse{
		EntitySourceType: info.Meta.EntitySourceType,
		Survivor:         res,
	}, nil
}
