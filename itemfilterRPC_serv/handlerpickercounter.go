package itemfilterRPC_serv

import (
	"context"

	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"
)

func (s *Server) GetPickerCounterNumber(ctx context.Context, in *itemfilterRPC_pb.GetPickerCounterNumberQuery) (*itemfilterRPC_pb.GetPickerCounterNumberResponse, error) {
	r, err := s.filterstore.GetPickerCounterNumber(in.Days, in.Offset, in.Mod, in.Ranges...)
	if err != nil {
		return nil, err
	}
	return &itemfilterRPC_pb.GetPickerCounterNumberResponse{Result: r}, nil
}
