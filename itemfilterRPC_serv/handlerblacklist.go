package itemfilterRPC_serv

import (
	"context"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetBlacklistList(ctx context.Context, in *emptypb.Empty) (*itemfilterRPC_pb.ListResponse, error) {
	r, err := s.metadatastore.GetBlacklistList()
	if err != nil {
		return nil, err
	}
	result := itemfilterRPC_pb.ListResponse{
		Content: r,
	}
	return &result, nil
}

// NewBlacklist 创建一个新的blacklist
func (s *Server) NewBlacklist(ctx context.Context, in *itemfilterRPC_pb.NewFilterQuery) (*itemfilterRPC_pb.IDResponse, error) {
	if in.Meta == nil {
		return nil, ErrNeedMetaInfo
	}
	if in.Meta.ID == "" {
		_blacklistid, err := s.idGener.Next()
		if err != nil {
			return nil, err
		}
		in.Meta.ID = _blacklistid
	} else {
		status, err := s.filterstore.BlacklistStatus(in.Meta.ID)
		if err != nil {
			return nil, err
		}
		if status.Exists {
			return nil, ErrIDExists
		}
	}
	in.Meta.Key = s.filterstore.BlacklistKey(in.Meta.ID)
	err := s.filterstore.NewBlacklist(in.Meta.ID, in.Setting, in.EntitySources...)
	if err != nil {
		return nil, err
	}
	err = s.metadatastore.NewBlacklist(in)
	if err != nil {
		//回退
		s.filterstore.DeleteBlacklist(in.Meta.ID)
		return nil, err
	}
	return &itemfilterRPC_pb.IDResponse{
		ID: in.Meta.ID,
	}, nil
}

// GetBlacklistInfo 获取黑名单信息
func (s *Server) GetBlacklistInfo(ctx context.Context, in *itemfilterRPC_pb.IDQuery) (*itemfilterRPC_pb.BlacklistInfoResponse, error) {
	if in.Id == "" {
		return nil, ErrNeedID
	}
	blacklistid := in.Id
	status, err := s.filterstore.BlacklistStatus(blacklistid)
	if err != nil {
		return nil, err
	}
	if !status.Exists {
		return nil, ErrNotFound
	}
	info, err := s.metadatastore.GetBlacklistInfo(blacklistid)
	if err != nil {
		return nil, err
	}
	info.Status = status
	return &itemfilterRPC_pb.BlacklistInfoResponse{Info: info}, nil
}

// deleteBlacklist 指定id后尝试删除这个id对应的信息
func (s *Server) deleteBlacklist(blacklistid string) error {
	err := s.filterstore.DeleteBlacklist(blacklistid)
	if err != nil {
		return err
	}
	err = s.metadatastore.DeleteBlacklist(blacklistid)
	if err != nil {
		return err
	}
	return nil
}

// DeleteBlacklist 删除blackList
func (s *Server) DeleteBlacklist(ctx context.Context, in *itemfilterRPC_pb.IDQuery) (*itemfilterRPC_pb.IDResponse, error) {
	if in.Id == "" {
		return nil, ErrNeedID
	}
	blacklistid := in.Id
	err := s.deleteBlacklist(blacklistid)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.IDResponse{
		ID: blacklistid,
	}, nil
}

// UpdateBlacklistSource 更新BlackList内不可用资源范围
func (s *Server) UpdateBlacklistSource(ctx context.Context, in *itemfilterRPC_pb.UpdateBlacklistSourceQuery) (*itemfilterRPC_pb.IDResponse, error) {
	if in.Id == "" {
		return nil, ErrNeedID
	}
	blacklistid := in.Id
	status, err := s.filterstore.BlacklistStatus(blacklistid)
	if err != nil {
		return nil, err
	}
	if !status.Exists {
		return nil, ErrNotFound
	}
	err = s.filterstore.UpdateBlacklistSource(blacklistid, in.Content...)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.IDResponse{ID: blacklistid}, nil
}

// CheckSourceInBlacklist 检查物品在不在黑名单中
func (s *Server) CheckSourceInBlacklist(ctx context.Context, in *itemfilterRPC_pb.CheckSourceQuery) (*itemfilterRPC_pb.CandidateStatusResponse, error) {
	if in.Id == "" {
		return nil, ErrNeedID
	}
	if len(in.Candidates) == 0 {
		return nil, ErrNeedAtLeastOneCandidate
	}

	blacklistid := in.Id
	status, err := s.filterstore.BlacklistStatus(blacklistid)
	if err != nil {
		return nil, err
	}
	if !status.Exists {
		return nil, ErrNotFound
	}
	info, err := s.metadatastore.GetBlacklistInfo(blacklistid)
	if err != nil {
		return nil, err
	}
	res, err := s.filterstore.CheckSourceInBlacklist(blacklistid, in.Candidates...)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.CandidateStatusResponse{
		EntitySourceType: info.Meta.EntitySourceType,
		CandidateStatus:  res,
	}, nil
}

// FilterSourceInBlacklist 过滤掉在黑名单中的物品,保留不在的部分
func (s *Server) FilterSourceInBlacklist(ctx context.Context, in *itemfilterRPC_pb.FilterSourceQuery) (*itemfilterRPC_pb.SurvivorResponse, error) {
	if in.Id == "" {
		return nil, ErrNeedID
	}
	if len(in.Candidates) == 0 {
		return nil, ErrNeedAtLeastOneCandidate
	}

	blacklistid := in.Id
	status, err := s.filterstore.BlacklistStatus(blacklistid)
	if err != nil {
		return nil, err
	}
	if !status.Exists {
		return nil, ErrNotFound
	}
	info, err := s.metadatastore.GetBlacklistInfo(blacklistid)
	if err != nil {
		return nil, err
	}
	res, err := s.filterstore.FilterSourceInBlacklist(blacklistid, in.Need, in.ChunkSize, in.Candidates...)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.SurvivorResponse{
		EntitySourceType: info.Meta.EntitySourceType,
		Survivor:         res,
	}, nil
}
