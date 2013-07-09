package models


type AvatarUser struct {
	UserId int              `json:"user_id"`
	UserName string         `json:"user_name"`
	UserPasswd string       `json:"user_passwd"`
	UserMobilePhone int     `json:"user_mobile_phone"`
	RoleType string			`json:"role_type"`
}

