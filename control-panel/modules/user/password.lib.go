package usermodule

import (
	"encoding/base64"
	"log"

	"golang.org/x/crypto/bcrypt"
)

const (
	COST = bcrypt.DefaultCost
)

func generateHash(password *string) (*string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashString := base64.StdEncoding.EncodeToString(hashBytes)

	return &hashString, nil
}

func verifyPassword(password *string, hash *string) bool {
	hashBytes, err := base64.StdEncoding.DecodeString(*hash)

	if err != nil {
		log.Println("[PasswordService] error while trying to decode hash: ", err.Error())
		return false
	}

	err = bcrypt.CompareHashAndPassword(hashBytes, []byte(*password))
	return err == nil
}
