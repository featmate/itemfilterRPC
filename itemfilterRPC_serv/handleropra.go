package itemfilterRPC_serv

import (
	"context"
	"io"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SetEntitySourceUsed  设置资源在指定条件中已被使用
func (s *Server) SetEntitySourceUsed(ctx context.Context, in *itemfilterRPC_pb.SetSourceUsedQuery) (*emptypb.Empty, error) {
	if len(in.EntitySources) <= 0 {
		return nil, ErrNeedAtLeastOneEntitySource
	}
	err := s.filterstore.SetEntitySourceUsed(in)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// BatchSetEntitySourceUsed  设置资源在指定条件中已被使用
func (s *Server) BatchSetEntitySourceUsed(stream itemfilterRPC_pb.ITEMFILTERRPC_BatchSetEntitySourceUsedServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&emptypb.Empty{})
			}
			return err
		}
		err = s.filterstore.SetEntitySourceUsed(in)
		if err != nil {
			return err
		}
	}
}

// CheckEntitySource 检查实体资源是否可用
func (s *Server) CheckEntitySource(ctx context.Context, in *itemfilterRPC_pb.CheckQuery) (*itemfilterRPC_pb.CandidateStatusResponse, error) {
	if in.Conditions == nil {
		return nil, ErrNeedConditions
	}
	if in.Conditions.EntitySourceType == "" {
		return nil, ErrNeedEntityType
	}
	if len(in.Candidates) <= 0 {
		return nil, ErrNeedAtLeastOneCandidate
	}
	res, err := s.filterstore.CheckEntitySource(in.Conditions, in.Candidates...)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.CandidateStatusResponse{EntitySourceType: in.Conditions.EntitySourceType, CandidateStatus: res}, nil
}

// FilterEntitySource 过滤掉不符合的,保留可用的
func (s *Server) FilterEntitySource(ctx context.Context, in *itemfilterRPC_pb.FilterQuery) (*itemfilterRPC_pb.SurvivorResponse, error) {
	if in.Conditions == nil {
		return nil, ErrNeedConditions
	}
	if in.Conditions.EntitySourceType == "" {
		return nil, ErrNeedEntityType
	}
	if len(in.Candidates) <= 0 {
		return nil, ErrNeedAtLeastOneCandidate
	}
	res, err := s.filterstore.FilterEntitySource(in.Conditions, in.Need, in.ChunkSize, in.Candidates...)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.SurvivorResponse{EntitySourceType: in.Conditions.EntitySourceType, Survivor: res}, nil
}
