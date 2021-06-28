package session_service

import (
	"auth2/mappers"
	"auth2/models"
	"auth2/storages"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type SessionsService struct {
	storage           *storages.SessionStorage
	sessionUserPrefix string
	sessionIdPrefix   string
	keyPartsSeparator string
	ttl               int
}

func New(storage *storages.SessionStorage, sessionUserPrefix string, sessionIdPrefix string, ttl int) *SessionsService {
	return &SessionsService{
		storage:           storage,
		sessionUserPrefix: sessionUserPrefix,
		sessionIdPrefix:   sessionIdPrefix,
		ttl:               ttl,
		keyPartsSeparator: ":",
	}
}

func (s *SessionsService) NewSession(userId int, sessionId string) *models.Session {
	var session models.Session
	session.UserId = userId
	session.Id = sessionId
	session.Created = (uint64)(time.Now().Unix())

	return &session
}

func (s *SessionsService) RemoveSession(userId int, sessionId string) {
	userSessionKey := s.GetUserSessionKey(userId, sessionId)
	sessionKey := s.GetSessionKey(sessionId)
	s.storage.DeleteKeys(userSessionKey, sessionKey)
}

func (s *SessionsService) RemoveSessions(userId int) {
	sessions := s.GetSessions(userId)
	for _, session := range sessions {
		s.RemoveSession(userId, session.Id)
	}
}

func (s *SessionsService) GetSessionBySessionId(sessionId string) *models.Session {
	pattern := s.sessionUserPrefix + "*:*" + sessionId
	keys, err := s.storage.SearchKeys(pattern)
	if err != nil || len(keys) < 1 {
		return nil
	}

	userId := s.ExtractUserIdFromUserSessionKey(keys[0])

	if userId == 0 {
		return nil
	}

	return s.GetSession(userId, sessionId)
}

func (s *SessionsService) GetSession(userId int, sessionId string) *models.Session {
	key := s.GetUserSessionKey(userId, sessionId)
	rawSession, err := s.storage.GetValue(key)
	if err != nil {
		return nil
	}

	session := mappers.MapRawSessionToSession(userId, s.ExtractSessionId(key), rawSession)

	return session
}

func (s *SessionsService) SaveSession(session *models.Session) {
	userSessionKey := s.GetUserSessionKey(session.UserId, session.Id)
	ttl, _ := time.ParseDuration(fmt.Sprintf("%ds", s.ttl))

	rawSession := mappers.MapSessionToRaw(session)

	s.storage.Save(userSessionKey, rawSession, ttl)
}

func (s *SessionsService) GetSessionTTL(userId int, sessionId string) time.Duration {
	key := s.GetUserSessionKey(userId, sessionId)
	duration, _ := s.storage.GetTTL(key)
	return duration
}

func (s *SessionsService) GetSessions(userId int) []models.Session {
	pattern := fmt.Sprintf("%s%s%s*", s.sessionUserPrefix, strconv.Itoa(userId), s.keyPartsSeparator)
	userSessionsKeys, err := s.storage.SearchKeys(pattern)

	if err != nil {
		logrus.Error("GetSessions error: " + err.Error())
	}

	sessions := make([]models.Session, len(userSessionsKeys))
	for i, key := range userSessionsKeys {
		sessionRawData, err := s.storage.GetValue(key)
		if err != nil {
			logrus.Error("Session data unavailable for key " + key + ": " + err.Error())
			continue
		}

		session := mappers.MapRawSessionToSession(userId, s.ExtractSessionId(key), sessionRawData)

		sessions[i] = *session
	}

	return sessions
}

func (s *SessionsService) GetUserSessionKey(userId int, sessionId string) string {
	return fmt.Sprintf("%s%s:%s",
		s.sessionUserPrefix,
		strconv.Itoa(userId),
		sessionId)
}

func (s *SessionsService) GetSessionKey(sessionId string) string {
	return fmt.Sprintf("%s%s",
		s.sessionIdPrefix,
		sessionId,
	)
}

func (s *SessionsService) ExtractSessionId(key string) string {
	// @todo remove duplicate
	sessionId := ""
	if strings.Index(key, s.sessionUserPrefix) != -1 {
		parts := strings.Split(key, s.keyPartsSeparator)
		if len(parts) > 0 {
			sessionId = parts[len(parts)-1]
		}

	} else if strings.Index(key, s.sessionIdPrefix) != -1 {
		parts := strings.Split(key, s.keyPartsSeparator)
		if len(parts) > 0 {
			sessionId = parts[len(parts)-1]
		}
	}

	return sessionId
}

func (s *SessionsService) ExtractUserIdFromUserSessionKey(key string) int {
	parts := strings.Split(key, s.keyPartsSeparator)
	if len(parts) < 2 {
		return 0
	}

	rawUserId := parts[len(parts)-2]

	userId, _ := strconv.Atoi(rawUserId)

	return userId
}
