package sessionmodule

import usermodule "platformlab/controlpanel/modules/user"

type Session struct {
	UserId uint
	Name   string
	Email  string
}

func NewSessionFromUser(user *usermodule.User) *Session {
	return &Session{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}
}
