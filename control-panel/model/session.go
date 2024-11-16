package model

type Session struct {
	UserId uint
	Name   string
	Email  string
}

func NewSessionFromUser(user *User) *Session {
	return &Session{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}
}
