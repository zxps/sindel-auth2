package server

import (
	pb "auth2/proto"
	"context"
)

func (s *AuthServer) SystemSupportNotify(ctx context.Context, in *pb.SystemSupportNotifyRequest) (*pb.EmptyResponse, error) {
	go func() {
		s.sendSupportNotify(in.GetMessage())
	}()

	return &pb.EmptyResponse{}, nil
}
