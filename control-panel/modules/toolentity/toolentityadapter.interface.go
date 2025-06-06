package toolentity

import tooleventmodule "synthreon/modules/toolevent"

type ToolEntityAdapter interface {
	StartHandler() error
	SendEvent(event *tooleventmodule.ToolEvent) error
	OnEventReceived(handler func(event *tooleventmodule.ToolEvent))
	OnDisconnect(handler func())
	Close()
}
