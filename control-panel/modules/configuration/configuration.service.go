package configurationmodule

import (
	"os"
	genericmodule "platformlab/controlpanel/modules/commonmodule"
)

type ConfigurationService struct {
	AccessTokenSecret string
	RootUserEmail     string
	RootPassword      string
	DatabasePath      string
}

func TryLoadApplicationConfigFromEnvironment() (*ConfigurationService, error) {
	secret := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if secret == "" {
		return nil, &genericmodule.GenericLogicError{Message: "[configuration] ACCESS_TOKEN_SECRET_KEY required"}
	}

	rootUserEmail := os.Getenv("ROOT_EMAIL")
	if rootUserEmail == "" {
		return nil, &genericmodule.GenericLogicError{Message: "[configuration] ROOT_EMAIL required"}
	}

	rootUserPassword := os.Getenv("ROOT_PASSWORD")
	if rootUserPassword == "" {
		return nil, &genericmodule.GenericLogicError{Message: "[configuration] ROOT_PASSWORD required"}
	}

	database := os.Getenv("DATABASE")
	if database == "" {
		return nil, &genericmodule.GenericLogicError{Message: "[configuration] DATABASE required"}
	}

	return &ConfigurationService{
		AccessTokenSecret: secret,
		RootUserEmail:     rootUserEmail,
		RootPassword:      rootUserPassword,
		DatabasePath:      database,
	}, nil
}
