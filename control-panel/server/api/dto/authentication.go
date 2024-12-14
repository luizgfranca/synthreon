package api

import sessionmodule "platformlab/controlpanel/modules/session"

type LoginRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginSuccessResopnseDto struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

func NewLoginSuccessResponseDto(session *sessionmodule.Session, accessToken *string) *LoginSuccessResopnseDto {
	return &LoginSuccessResopnseDto{
		ID:          session.UserId,
		Name:        session.Name,
		Email:       session.Email,
		AccessToken: *accessToken,
	}
}
