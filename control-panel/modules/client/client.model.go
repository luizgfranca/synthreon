package clientmodule

import (
	"log"
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
	usermodule "platformlab/controlpanel/modules/user"
)

type Manager interface {
	FindProject(acronym string) (*projectmodule.Project, error)
	FindTool(project *projectmodule.Project, acronym string) (*toolmodule.Tool, error)

	DistributeEvent(e *tooleventmodule.ToolEvent)
}

type Client struct {
	ID        string
	manager   Manager
	user      *usermodule.User
	entity    toolentity.ToolEntityAdapter
	terminals []Terminal
}

func (c *Client) log(v ...any) {
	x := append([]any{"[Provider-" + c.ID + "]"}, v...)

	log.Println(x...)
}
