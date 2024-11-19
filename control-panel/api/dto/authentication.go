package api

import "platformlab/controlpanel/model"

type LoginRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginSuccessResopnseDto struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewLoginSuccessResponseDto(user *model.User) *LoginSuccessResopnseDto {
	return &LoginSuccessResopnseDto{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
