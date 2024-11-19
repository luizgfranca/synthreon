package service

import (
	"fmt"
	"platformlab/controlpanel/model"
	"platformlab/controlpanel/util"

	"gorm.io/gorm"
)

type User struct {
	Db *gorm.DB
}

func (u *User) FindByEmail(email string) (*model.User, error) {
	var maybeUser *model.User

	result := u.Db.Where("email = ?", email).First(&maybeUser)
	if result.Error != nil {
		return nil, &model.GenericLogicError{
			Message: fmt.Sprintf("user with email %s not found", email),
		}
	}

	return maybeUser, nil
}

func (u *User) VerifyAuthenticationCredentials(email *string, password *string) (*model.User, error) {
	user, err := u.FindByEmail(*email)
	if err != nil {
		return nil, err
	}

	isValid := util.VerifyPassword(password, &user.Hash)
	if !isValid {
		return nil, &model.GenericLogicError{
			Message: fmt.Sprintf("invalid password for user %s", *email),
		}
	}

	return user, nil
}

func (u *User) Create(user *model.User) (*model.User, error) {
	var result *gorm.DB

	_, err := u.FindByEmail(user.Email)
	if err == nil {
		return nil, &model.GenericLogicError{
			Message: fmt.Sprintf("element with email %s already exists", user.Email),
		}
	}

	result = u.Db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	created, err := u.FindByEmail(user.Email)
	if err != nil {
		return nil, result.Error
	}
	if created == nil {
		panic("created item in database, but it was not found after insertion")
	}

	return created, nil
}
