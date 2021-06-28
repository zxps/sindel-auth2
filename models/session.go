package models

type Session struct {
	Id        string `json:"-"`
	UserId    int    `json:"-"`
	LastUri   string `json:"last_uri"`
	UserAgent string `json:"user_agent"`
	Ip        uint64 `json:"ip"`
	Updated   uint64 `json:"updated"`
	Created   uint64 `json:"created"`
}
