package server

import (
	"auth2/mappers"
	"auth2/models"
	pb "auth2/proto"
	"context"
	"fmt"
	"time"
)

func (s *AuthServer) DeleteUserSession(ctx context.Context, in *pb.UserSessionRequest) (*pb.EmptyResponse, error) {
	session, err := s.getSessionBySessionRequest(in)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	sessions := s.getSessionService()

	sessions.RemoveSession(session.UserId, session.Id)

	return &pb.EmptyResponse{}, nil
}

func (s *AuthServer) SaveUserSession(ctx context.Context, in *pb.UserSessionRequest) (*pb.EmptyResponse, error) {
	session, err := s.getSessionBySessionRequest(in)
	if err != nil {
		return nil, err
	}

	sessionsService := s.getSessionService()
	if session == nil {
		session = sessionsService.NewSession((int)(in.GetUserId()), in.GetSessionId())
	}

	if in.UserAgent != nil {
		session.UserAgent = in.GetUserAgent()
	}
	if in.Ip != nil {
		session.Ip = in.GetIp()
	}

	if in.LastUri != nil {
		session.LastUri = in.GetLastUri()
	}

	session.Updated = (uint64)(time.Now().Unix())

	sessionsService.SaveSession(session)

	return &pb.EmptyResponse{}, nil
}

func (s *AuthServer) GetUserSessions(ctx context.Context, in *pb.UserRequest) (*pb.GetUserSessionsResponse, error) {
	user, err := s.getUserByUserRequestProto(in)
	if err != nil {
		return &pb.GetUserSessionsResponse{}, err
	}

	sessionsService := s.getSessionService()
	userSessions := sessionsService.GetSessions(user.Id)
	var sessions = make([]*pb.UserSession, len(userSessions))
	for i, session := range userSessions {
		proto := mappers.MapSessionToProto(&session)

		ttl := sessionsService.GetSessionTTL(session.UserId, session.Id)
		proto.Ttl = (uint64)(ttl.Seconds())
		proto.TtlString = ttl.String()

		sessions[i] = proto
	}
	return &pb.GetUserSessionsResponse{
		Sessions: sessions,
	}, nil
}

func (s *AuthServer) getSessionBySessionRequest(in *pb.UserSessionRequest) (*models.Session, error) {
	var session *models.Session
	if in.UserId != nil {
		if in.GetUserId() == 0 {
			return nil, fmt.Errorf("wrong user id")
		}
		usersService := s.getUserService()
		if !usersService.IsExists((int)(in.GetUserId())) {
			return nil, fmt.Errorf("user not exists")
		}
	}

	sessions := s.getSessionService()

	if in.UserId == nil {
		session = sessions.GetSessionBySessionId(in.GetSessionId())
	} else {
		session = sessions.GetSession((int)(in.GetUserId()), in.GetSessionId())
	}

	if session == nil && in.UserId == nil {
		return nil, fmt.Errorf("session not found")
	}

	return session, nil
}
