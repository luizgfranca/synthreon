package util

import (
	"os"
)

type Configuration struct {
	AccessTokenSecret string
	RootUserEmail     string
	RootPassword      string
	DatabasePath      string
}

func TryLoadApplicationConfigFromEnvironment() (*Configuration, error) {
	secret := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if secret == "" {
		return nil, &GenericLogicError{Message: "[configuration] ACCESS_TOKEN_SECRET_KEY required"}
	}

	rootUserEmail := os.Getenv("ROOT_EMAIL")
	if rootUserEmail == "" {
		return nil, &GenericLogicError{Message: "[configuration] ROOT_EMAIL required"}
	}

	rootUserPassword := os.Getenv("ROOT_PASSWORD")
	if rootUserPassword == "" {
		return nil, &GenericLogicError{Message: "[configuration] ROOT_PASSWORD required"}
	}

	database := os.Getenv("DATABASE")
	if database == "" {
		return nil, &GenericLogicError{Message: "[configuration] DATABASE required"}
	}

	return &Configuration{
		AccessTokenSecret: secret,
		RootUserEmail:     rootUserEmail,
		RootPassword:      rootUserPassword,
		DatabasePath:      database,
	}, nil
}
