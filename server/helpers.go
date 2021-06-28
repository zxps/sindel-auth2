package server

import (
	"auth2/models"
	pb "auth2/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) getUserByUserRequestProto(proto *pb.UserRequest) (*models.User, error) {
	if !s.isValidUserRequestProto(proto) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user request")
	}

	userService := s.getUserService()

	var user *models.User
	var err error
	if proto.UserId != nil {
		user, err = userService.GetById(int(proto.GetUserId()))
	} else if proto.Username != nil {
		user, err = userService.GetByUsername(proto.GetUsername())
	} else if proto.Email != nil {
		user, err = userService.GetByEmail(proto.GetEmail())
	}

	var protoError error
	if err != nil {
		protoError = status.Errorf(codes.NotFound, "user not found (%s)", err.Error())
	}

	return user, protoError
}

func (s *AuthServer) isValidUserRequestProto(proto *pb.UserRequest) bool {
	if proto.UserId != nil {
		return true
	}

	if proto.Username != nil {
		return true
	}

	if proto.Email != nil {
		return true
	}

	return false
}
