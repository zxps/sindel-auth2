package mappers

import (
	"auth2/models"
	pb "auth2/proto"
	"auth2/utils"
	"encoding/json"
)

func MapSessionToProto(s *models.Session) *pb.UserSession {
	var session pb.UserSession

	session.Id = s.Id
	session.LastUri = s.LastUri
	session.Ip = s.Ip
	session.IpString = utils.ConvertIpToString(s.Ip)
	session.Created = s.Created
	session.Updated = s.Updated
	session.UserAgent = s.UserAgent

	return &session
}

func MapRawSessionToSession(userId int, sessionId string, rawSession string) *models.Session {
	var session models.Session

	json.Unmarshal([]byte(rawSession), &session)
	session.Id = sessionId
	session.UserId = userId

	return &session
}

func MapSessionToRaw(s *models.Session) string {
	result, _ := json.Marshal(s)

	return (string)(result)
}
