package sessionmodule

import (
	usermodule "platformlab/controlpanel/modules/user"

	"github.com/google/uuid"
)

type Session struct {
	ID     string
	UserId uint
	Name   string
	Email  string
}

func NewSessionFromUser(user *usermodule.User) *Session {
	return &Session{
		ID:     uuid.NewString(),
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}
}
