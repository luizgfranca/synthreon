package usermodule

import (
	"fmt"

	commonmodule "platformlab/controlpanel/modules/common"

	"gorm.io/gorm"
)

type UserService struct {
	Db *gorm.DB
}

func (u *UserService) FindByEmail(email string) (*User, error) {
	var maybeUser *User

	result := u.Db.Where("email = ?", email).First(&maybeUser)
	if result.Error != nil {
		return nil, &commonmodule.GenericLogicError{
			Message: fmt.Sprintf("user with email %s not found", email),
		}
	}

	return maybeUser, nil
}

func (u *UserService) VerifyAuthenticationCredentials(email *string, password *string) (*User, error) {
	user, err := u.FindByEmail(*email)
	if err != nil {
		return nil, err
	}

	isValid := verifyPassword(password, &user.Hash)
	if !isValid {
		return nil, &commonmodule.GenericLogicError{
			Message: fmt.Sprintf("invalid password for user %s", *email),
		}
	}

	return user, nil
}

func (u *UserService) Create(user *User) (*User, error) {
	var result *gorm.DB

	_, err := u.FindByEmail(user.Email)
	if err == nil {
		return nil, &commonmodule.GenericLogicError{
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
