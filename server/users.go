package server

import (
	"auth2/mappers"
	"auth2/models"
	pb "auth2/proto"
	"auth2/storages"
	"auth2/validator"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (s *AuthServer) CreateUser(ctx context.Context, in *pb.UserPayload) (*pb.User, error) {
	if len(in.Email) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "email required")
	}

	if err := validator.EmailValidate(in.Email); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email (%s)", err.Error())
	}

	in.Email = validator.EmailNormalize(in.Email)

	if len(in.Username) < 1 {
		in.Username = in.Email
	}

	userService := s.getUserService()
	security := s.getSecurityService()

	if userService.IsExistsByEmail(in.GetEmail()) {
		return nil, status.Errorf(codes.AlreadyExists, "user with email '%s' already exists", in.GetEmail())
	}

	if userService.IsExistsByUsername(in.GetUsername()) {
		return nil, status.Errorf(codes.AlreadyExists, "user with username '%s' already exists", in.GetUsername())
	}

	if in.Password != nil {
		passwdValidateOptions := &validator.PasswordValidateOptions{
			MinLength: s.container.Config().UserPasswordMinLength,
			MaxLength: s.container.Config().UserPasswordMaxLength,
		}
		if err := validator.PasswordValidate(in.GetPassword(), passwdValidateOptions); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		password, err := security.HashPasswordAsString(in.GetPassword())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		in.Password = &password
	}

	if in.Id != nil {
		in.Id = nil
	}

	in.IsModerated = false
	in.IsEnabled = false

	var user models.User
	mappers.MapProtoUserPayloadToUser(in, &user)

	if in.RoleIds != nil {
		roles := userService.GetRoles(in.GetRoleIds())
		user.SetRoles(roles)
	}

	err := userService.CreateUser(&user)
	if err != nil || user.Id <= 0 {
		return &pb.User{}, status.Error(codes.Internal, err.Error())
	}

	confirmationToken := s.getTokenService().GenerateUserConfirmation(&user)
	if len(confirmationToken) <= 0 {
		userService.DeleteUser(&user)
		return nil, status.Errorf(codes.Internal, "unable to generate confirmation token")
	}

	secretUrl := strings.Replace(s.getConfig().UserConfirmUrlPattern, "{token}", confirmationToken, -1)

	go func() {
		s.getMailService().SendRegistrationConfirmation(user.Email, secretUrl)
		s.sendSupportNotify(fmt.Sprintf("Зарегистрировался пользователь %s (%s)", user.Username, user.Email))
	}()

	var protoUser pb.User

	mappers.MapUserToUserProto(&user, &protoUser)

	return &protoUser, nil
}

func (s *AuthServer) UpdateUser(ctx context.Context, in *pb.UserPayload) (*pb.User, error) {
	userService := s.getUserService()

	userId := (int)(in.GetId())
	if !userService.IsExists(userId) {
		return nil, fmt.Errorf("user not found")
	}

	user, err := userService.GetById(userId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if len(in.Email) < 1 {
		in.Email = user.Email
	}

	if len(in.Username) < 1 {
		in.Username = user.Username
	}

	if in.Email != user.Email {
		if err := validator.EmailValidate(in.Email); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid email", in.Email)
		}

		if userService.IsExistsByEmail(in.Email) {
			return nil, status.Errorf(codes.InvalidArgument, "user with email '%s' already exists", in.Email)
		}
	}

	if in.Username != user.Username {
		if userService.IsExistsByUsername(in.Username) {
			return nil, status.Errorf(codes.InvalidArgument, "user with username '%s' already exists", in.Username)
		}
	}

	if in.Password != nil && in.GetPassword() != user.Password {
		return nil, status.Errorf(codes.InvalidArgument, "change of password is not available")
	}

	if in.RoleIds != nil {
		roles := userService.GetRoles(in.GetRoleIds())
		user.SetRoles(roles)
	}

	mappers.MapProtoUserPayloadToUser(in, user)

	if err := userService.UpdateUser(user); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var proto pb.User
	mappers.MapUserToUserProto(user, &proto)

	return &proto, err
}

func (s *AuthServer) SearchUsers(ctx context.Context, in *pb.SearchUsersRequest) (*pb.UsersResponse, error) {
	userService := s.getUserService()

	scanOptions := &storages.ScanOptions{
		Limit:  10,
		Offset: 0,
	}

	if in.Offset != nil {
		scanOptions.Offset = uint(in.GetOffset())
	}

	if in.Limit != nil {
		scanOptions.Limit = uint(in.GetLimit())
	}

	scanOptions.IsDesc = in.IsDesc

	users, err := userService.Scan(scanOptions)

	if err != nil {
		return &pb.UsersResponse{}, err
	}

	var usersResponse pb.UsersResponse
	usersResponse.Users = make([]*pb.User, len(users))

	for i, user := range users {
		var proto pb.User
		mappers.MapUserToUserProto(&user, &proto)
		usersResponse.Users[i] = &proto
	}

	if in.GetIsCountTotal() {
		totalCount := userService.GetTotalCount()
		convertedTotalCount := int32(totalCount)
		usersResponse.TotalCount = &convertedTotalCount
	}

	return &usersResponse, nil
}

func (s *AuthServer) GetUser(ctx context.Context, in *pb.UserRequest) (*pb.User, error) {
	user, err := s.getUserByUserRequestProto(in)

	if err != nil {
		return &pb.User{}, err
	}

	var proto pb.User

	mappers.MapUserToUserProto(user, &proto)

	return &proto, nil
}

func (s *AuthServer) DeleteUser(ctx context.Context, in *pb.UserRequest) (*pb.EmptyResponse, error) {
	user, err := s.getUserByUserRequestProto(in)

	if err != nil {
		return &pb.EmptyResponse{}, err
	}

	err = s.getUserService().DeleteUser(user)

	if err != nil {
		return &pb.EmptyResponse{}, nil
	}

	s.getSessionService().RemoveSessions(user.Id)

	return &pb.EmptyResponse{}, nil
}

func (s *AuthServer) DeleteUsers(ctx context.Context, in *pb.UserIdsRequest) (*pb.EmptyResponse, error) {
	var users []*models.User = make([]*models.User, len(in.GetUserIds()))

	// Check user ids first
	for i, userId := range in.GetUserIds() {
		user, err := s.getUserService().GetById(int(userId))
		if err != nil {
			return nil, err
		}
		users[i] = user
	}

	// Remove
	for _, user := range users {
		s.getUserService().DeleteUser(user)
	}

	return &pb.EmptyResponse{}, nil
}

func (s *AuthServer) SetUserPassword(ctx context.Context, in *pb.UserPasswordRequest) (*pb.EmptyResponse, error) {
	user, err := s.getUserByUserRequestProto(in.GetUser())

	if err != nil {
		return &pb.EmptyResponse{}, err
	}

	userService := s.getUserService()

	hashedPassword, err := s.getSecurityService().HashPassword(in.GetPassword())
	if err != nil {
		return &pb.EmptyResponse{}, status.Errorf(codes.Internal, "unable to create password hash (%s)", err.Error())
	}

	err = userService.UpdatePassword(user.Username, hashedPassword)
	if err != nil {
		return &pb.EmptyResponse{}, status.Errorf(codes.Internal, "unable to update user with new password")
	}

	return &pb.EmptyResponse{}, nil
}
