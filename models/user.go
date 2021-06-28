package models

import (
	"database/sql"
	"encoding/json"
)

type User struct {
	Id                  int        `db:"id" json:"id" goqu:"skipinsert"`
	Username            string     `db:"username" json:"username"`
	ImageId             NullInt32  `db:"image_id" json:"image_id"`
	FirstName           string     `db:"first_name" json:"first_name"`
	LastName            string     `db:"last_name" json:"last_name"`
	Email               string     `db:"email" json:"email"`
	Phone               string     `db:"phone" json:"phone"`
	Password            string     `db:"password" json:"-"`
	IsEnabled           bool       `db:"is_enabled" json:"is_enabled"`
	IsModerated         bool       `db:"is_moderated" json:"is_moderated"`
	LastLogin           NullString `db:"last_login" json:"last_login"`
	LastIP              uint64     `db:"last_ip" json:"last_ip"`
	ConfirmationToken   NullString `db:"confirmation_token" json:"-"`
	PasswordRequestedAt NullString `db:"password_requested_at" json:"-"`
	Created             string     `db:"created" json:"created"`
	Updated             string     `db:"updated" json:"updated"`
	Roles               UserRoles  `db:"roles" json:"roles"`
	Params              Params     `db:"params" json:"params"`
	Salt                string     `db:"salt" json:"-"`
}

func (u *User) SetRoles(roles []*Role) {
	if roles == nil || len(roles) < 1 {
		u.Roles.Scan("[]")
		return
	}

	uniqueRoles := map[string]*Role{}

	for _, role := range roles {
		if _, exists := uniqueRoles[role.Id]; !exists {
			uniqueRoles[role.Id] = role
		}
	}

	var rolesArray []*Role
	for _, role := range uniqueRoles {
		rolesArray = append(rolesArray, role)
	}

	rolesJSON, err := json.Marshal(rolesArray)
	if err == nil {
		u.Roles.Scan(rolesJSON)
	} else {
		u.Roles.Scan("[]")
	}
}

// UserRoles struct
type UserRoles struct {
	sql.NullString
}

// MarshalJSON for roles
func (r UserRoles) MarshalJSON() ([]byte, error) {
	if r.Valid && len(r.String) > 0 {
		var result []map[string]interface{}

		json.Unmarshal([]byte(r.String), &result)
		return json.Marshal(result)
	}

	return json.Marshal([]map[string]interface{}{})
}
