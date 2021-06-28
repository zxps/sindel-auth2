package mappers

import (
	"auth2/models"
	pb "auth2/proto"
	"auth2/utils"
)

func MapProtoUserPayloadToUser(proto *pb.UserPayload, u *models.User) {
	if proto.Id != nil {
		u.Id = (int)(proto.GetId())
	}

	u.Username = proto.GetUsername()
	u.Email = proto.GetEmail()
	u.FirstName = proto.GetFirstName()
	u.LastName = proto.GetLastName()
	u.LastIP = uint64(proto.LastIp)
	u.IsEnabled = proto.GetIsEnabled()
	u.IsModerated = proto.GetIsModerated()
	u.Phone = proto.GetPhone()

	if proto.LastLogin != nil {
		u.LastLogin.Scan(utils.ConvertTimestampToDatetime((int64)(proto.GetLastLogin())))
	}

	if proto.Password != nil {
		u.Password = proto.GetPassword()
	}

	if proto.ImageId != nil {
		u.ImageId.Scan(proto.GetImageId())
	}

	if proto.PasswordRequestedAt != nil {
		u.PasswordRequestedAt.Scan(utils.ConvertTimestampToDatetime((int64)(proto.GetPasswordRequestedAt())))
	}
	if proto.Params != nil {
		u.Params.Scan(proto.GetParams())
	}
	if proto.Roles != nil {
		u.Roles.Scan(proto.GetRoles())
	}

	if proto.Created != nil {
		u.Created = utils.ConvertTimestampToDatetime((int64)(proto.GetCreated()))
	}

	if proto.Updated != nil {
		u.Updated = utils.ConvertTimestampToDatetime((int64)(proto.GetUpdated()))
	}
}

func MapProtoUserToUser(proto *pb.User, u *models.User) {
	u.Username = proto.GetUsername()
	u.Email = proto.GetEmail()
	u.FirstName = proto.GetFirstName()
	u.LastName = proto.GetLastName()
	u.LastIP = uint64(proto.LastIp)
	u.IsEnabled = proto.GetIsEnabled()
	u.IsModerated = proto.GetIsModerated()
	u.Phone = proto.GetPhone()
	u.Password = proto.GetPassword()

	if proto.ImageId != nil {
		u.ImageId.Scan(proto.GetImageId())
	}

	if proto.PasswordRequestedAt != nil {
		u.PasswordRequestedAt.Scan(proto.GetPasswordRequestedAt())
	}
	if proto.Params != nil {
		u.Params.Scan(proto.GetParams())
	}
	if proto.Roles != nil {
		u.Roles.Scan(proto.GetRoles())
	}

	u.Created = utils.ConvertTimestampToDatetime((int64)(proto.GetCreated()))
	u.Updated = utils.ConvertTimestampToDatetime((int64)(proto.GetUpdated()))
}

func MapUserToUserProto(user *models.User, proto *pb.User) {
	proto.Id = uint32(user.Id)
	proto.Username = user.Username
	proto.Email = user.Email
	proto.Phone = user.Phone
	proto.LastIp = uint32(user.LastIP)
	proto.FirstName = user.FirstName
	proto.LastName = user.LastName
	proto.Password = user.Password
	proto.IsEnabled = user.IsEnabled
	proto.IsModerated = user.IsModerated

	if user.Roles.Valid {
		proto.Roles = &user.Roles.String
	}

	if user.Params.Valid {
		proto.Params = &user.Params.String
	}

	if user.ImageId.Valid {
		imageId := uint32(user.ImageId.Int32)
		proto.ImageId = &imageId
	}

	if user.LastLogin.Valid {
		ts := utils.ConvertDatetimeToTimestamp(user.LastLogin.String)
		proto.LastLogin = &ts
	}
	if user.PasswordRequestedAt.Valid {
		ts := utils.ConvertDatetimeToTimestamp(user.PasswordRequestedAt.String)
		proto.PasswordRequestedAt = &ts
	}

	proto.Created = utils.ConvertDatetimeToTimestamp(user.Created)
	proto.Updated = utils.ConvertDatetimeToTimestamp(user.Updated)
}
