package configurationmodule

import (
	"os"
	commonmodule "platformlab/controlpanel/modules/common"
	"strconv"
	"strings"
)

type ConfigurationService struct {
	AccessTokenSecret   string
	RootUserEmail       string
	RootPassword        string
	DatabasePath        string
	StaticFilesDir      string
	RetryTimeoutSeconds int

	// optional parameters
	AllowToolAutoCreation bool
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

	retryTimeoutStr := os.Getenv("RETRY_TIMEOUT_SECONDS")
	if retryTimeoutStr == "" {
		return nil, &commonmodule.GenericLogicError{Message: "[configuration] RETRY_TIMEOUT_SECONDS required"}
	}
	retryTimeoutSeconds, err := strconv.Atoi(retryTimeoutStr)
	if err != nil {
		return nil, &commonmodule.GenericLogicError{Message: "[configuration] RETRY_TIMEOUT_SECONDS should be an integer"}
	}

	allowToolAutoCreation := true
	allowToolAutoCreationStr := os.Getenv("ALLOW_TOOL_AUTOCREATION")
	if strings.ToUpper(allowToolAutoCreationStr) == "TRUE" {
		allowToolAutoCreation = true
	} else if strings.ToUpper(allowToolAutoCreationStr) == "FALSE" {
		allowToolAutoCreation = false
	}

	return &ConfigurationService{
		AccessTokenSecret:     secret,
		RootUserEmail:         rootUserEmail,
		RootPassword:          rootUserPassword,
		DatabasePath:          database,
		StaticFilesDir:        staticFilesDir,
		RetryTimeoutSeconds:   retryTimeoutSeconds,
		AllowToolAutoCreation: allowToolAutoCreation,
	}, nil
}
