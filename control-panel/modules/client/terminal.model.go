package clientmodule

import (
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
)

type Terminal struct {
	ID      string
	Project *projectmodule.Project
	Tool    *toolmodule.Tool
}
