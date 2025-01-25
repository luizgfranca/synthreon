package providermodule

import (
	"log"
	commonmodule "platformlab/controlpanel/modules/common"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
)

// TODO: should i make this data serailizeable? (for now it just expects events to arrive ordered already)
type ContextProviderResolver struct {
	contextToProviderAssignment map[string]*Provider
}

// Register
// provider should already be registered
func (c *ContextProviderResolver) Register(contextId string, provider *Provider) {
	if c.contextToProviderAssignment == nil {
		c.contextToProviderAssignment = map[string]*Provider{}
	}

	log.Println("[ContextProviderResolver] registering context: ", contextId, "to provider: ", provider.ID)
	if _, ok := c.contextToProviderAssignment[contextId]; ok {
		panic("unexpected error: tryiong to register a provider to handle an already assigned context")
	}

	c.contextToProviderAssignment[contextId] = provider
}

func (c *ContextProviderResolver) Unregister(contextId string) {
	log.Println("[ContextProviderResolver] unregistering context: ", contextId)
	delete(c.contextToProviderAssignment, contextId)
}

func (c *ContextProviderResolver) TryRouteEvent(e *tooleventmodule.ToolEvent) error {
	if c.contextToProviderAssignment == nil {
		c.contextToProviderAssignment = map[string]*Provider{}
	}

	contextId := e.ContextId

	log.Println("[ContextProviderResolver] trying to route event from context: ", e.ContextId)
	if _, ok := c.contextToProviderAssignment[contextId]; !ok {
		return &ContextNotFounError{}
	}

	success := c.contextToProviderAssignment[contextId].SendEvent(e)
	if !success {
		return &commonmodule.GenericLogicError{Message: "error sending event"}
	}

	return nil
}

// TryResolve return is nullable
func (c *ContextProviderResolver) TryResolve(ctxid string) *Provider {
	provider, ok := c.contextToProviderAssignment[ctxid]
	if !ok {
		return nil
	}

	return provider
}

type ContextNotFounError struct{}

func (c *ContextNotFounError) Error() string {
	return "context not found"
}
