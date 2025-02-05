package configurationmodule

import (
	"os"
	commonmodule "platformlab/controlpanel/modules/common"
)

type ConfigurationService struct {
	AccessTokenSecret string
	RootUserEmail     string
	RootPassword      string
	DatabasePath      string
	StaticFilesDir    string
}

func TryLoadApplicationConfigFromEnvironment() (*ConfigurationService, error) {
	secret := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if secret == "" {
		return nil, &commonmodule.GenericLogicError{Message: "[configuration] ACCESS_TOKEN_SECRET_KEY required"}
	}

	rootUserEmail := os.Getenv("ROOT_EMAIL")
	if rootUserEmail == "" {
		return nil, &commonmodule.GenericLogicError{Message: "[configuration] ROOT_EMAIL required"}
	}

	rootUserPassword := os.Getenv("ROOT_PASSWORD")
	if rootUserPassword == "" {
		return nil, &commonmodule.GenericLogicError{Message: "[configuration] ROOT_PASSWORD required"}
	}

	database := os.Getenv("DATABASE")
	if database == "" {
		return nil, &commonmodule.GenericLogicError{Message: "[configuration] DATABASE required"}
	}

	staticFilesDir := os.Getenv("STATIC_FILES_DIR")
	if staticFilesDir == "" {
		return nil, &commonmodule.GenericLogicError{Message: "[configuration] STATIC_FILES_DIR required"}
	}

	return &ConfigurationService{
		AccessTokenSecret: secret,
		RootUserEmail:     rootUserEmail,
		RootPassword:      rootUserPassword,
		DatabasePath:      database,
		StaticFilesDir:    staticFilesDir,
	}, nil
}
