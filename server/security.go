package server

import (
	"auth2/mappers"
	pb "auth2/proto"
	"auth2/validator"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (s *AuthServer) SetUserRoles(ctx context.Context, in *pb.SetUserRolesRequest) (*pb.User, error) {
	user, err := s.getUserByUserRequestProto(in.GetUser())
	if err != nil {
		return nil, err
	}

	err = validator.ValidateUserRoleIds(in.GetRoleIds(), s.getUserService())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	userService := s.getUserService()

	user.SetRoles(userService.GetRoles(in.GetRoleIds()))

	userService.UpdateUser(user)

	var userProto pb.User

	mappers.MapUserToUserProto(user, &userProto)

	return &userProto, nil
}

func (s *AuthServer) RequestUserPassword(ctx context.Context, in *pb.UserRequest) (*pb.User, error) {
	user, err := s.getUserByUserRequestProto(in)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty arguments or user not exists")
	}

	if !user.IsEnabled {
		return nil, status.Errorf(UserNotEnabledStatusCode, "user not confirmed/enabled")
	}

	if !user.IsModerated {
		return nil, status.Errorf(UserNotModeratedStatusCode, "user not moderated")
	}

	requestPasswordToken := s.getTokenService().GeneratePasswordRequest(user)

	var proto pb.User
	mappers.MapUserToUserProto(user, &proto)

	secretUrl := strings.Replace(s.getConfig().UserConfirmPasswordUrlPattern, "{token}", requestPasswordToken, -1)

	go func() {
		s.getMailService().SendPasswordRequest(user.Email, secretUrl)
		s.sendSupportNotify(fmt.Sprintf("Пользователь %s (%s) запросил новый пароль", user.Username, user.Email))
	}()

	return &proto, nil
}

func (s *AuthServer) ChangeUserPassword(ctx context.Context, in *pb.ChangeUserPasswordRequest) (*pb.User, error) {
	if len(in.Token) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "token empty or invalid")
	}

	ok, userId := s.getTokenService().FindPasswordRequest(in.GetToken())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "token not found")
	}

	user, err := s.getUserService().GetById(userId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if !user.IsEnabled {
		return nil, status.Errorf(UserNotEnabledStatusCode, "user not enabled/confirmed")
	}

	if !user.IsModerated {
		return nil, status.Errorf(UserNotModeratedStatusCode, "user not moderated")
	}

	if err := validator.PasswordValidate(in.GetPassword(), &validator.PasswordValidateOptions{
		MinLength: s.getConfig().UserPasswordMinLength,
		MaxLength: s.getConfig().UserPasswordMaxLength,
	}); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	passwordHash, err := s.getSecurityService().HashPasswordAsString(in.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to create password hash (%s)", err.Error())
	}

	user.Password = passwordHash

	err = s.getUserService().UpdateUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to save user (%s)", err.Error())
	}

	s.getTokenService().DeletePasswordRequest(in.GetToken())

	var protoUser pb.User

	mappers.MapUserToUserProto(user, &protoUser)

	return &protoUser, nil
}

func (s *AuthServer) CheckUserPassword(ctx context.Context, in *pb.UserPasswordRequest) (*pb.EmptyResponse, error) {
	user, err := s.getUserByUserRequestProto(in.GetUser())
	if err != nil {
		return nil, err
	}
	securityService := s.getSecurityService()
	err = securityService.CompareStringPasswords(user.Password, in.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("wrong password")
	}

	return &pb.EmptyResponse{}, nil
}

func (s *AuthServer) ConfirmUser(ctx context.Context, in *pb.TokenRequest) (*pb.User, error) {
	if len(in.GetToken()) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "token not specified")
	}

	ok, userId := s.getTokenService().FindUserConfirmation(in.GetToken())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "token not found")
	}

	user, err := s.getUserService().GetById(userId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if user.IsEnabled {
		return nil, status.Errorf(UserAlreadyConfirmed, "user already confirmed")
	}

	user.IsEnabled = true

	err = s.getUserService().UpdateUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to update user (%s)", err.Error())
	}

	s.getTokenService().DeleteUserConfirmation(in.GetToken())

	go func() {
		s.getMailService().SendUserConfirmed(user.Email)
	}()

	var proto pb.User
	mappers.MapUserToUserProto(user, &proto)

	return &proto, nil
}
