package clientmodule

import (
	projectmodule "synthreon/modules/project"
	toolmodule "synthreon/modules/tool"
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
