package providermodule

import (
	toolmodule "synthreon/modules/tool"

	"github.com/google/uuid"
)

type HandlerStatus string

const (
	HandlerStatusRegistering HandlerStatus = "REGISTERING"
	HandlerStatusActive      HandlerStatus = "ACTIVE"
	HandlerStatusInactive    HandlerStatus = "INACTIVE"
)

type Handler struct {
	ID     string
	Tool   *toolmodule.Tool
	Status HandlerStatus
}

func NewHandler(t *toolmodule.Tool) Handler {
	id := uuid.NewString()

	return Handler{
		ID:     id,
		Tool:   t,
		Status: HandlerStatusActive,
	}
}
