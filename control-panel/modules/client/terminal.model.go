package clientmodule

import (
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
)

type TerminalStatus string

const (
	TerminalStatusRunning  TerminalStatus = "RUNNING"
	TerminalStatusFinished TerminalStatus = "FINISHED"
)

type Terminal struct {
	ID      string
	Client  *Client
	Status  TerminalStatus
	Project *projectmodule.Project
	Tool    *toolmodule.Tool
}
