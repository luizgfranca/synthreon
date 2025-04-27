package providermodule

import (
	"log"
	commonmodule "synthreon/modules/common"
	tooleventmodule "synthreon/modules/toolevent"
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

	c.log("registering context: ", contextId, "to provider: ", provider.ID)
	if _, ok := c.contextToProviderAssignment[contextId]; ok {
		panic("unexpected error: tryiong to register a provider to handle an already assigned context")
	}

	c.contextToProviderAssignment[contextId] = provider
}

func (c *ContextProviderResolver) Unregister(contextId string) {
	c.log("unregistering context: ", contextId)
	delete(c.contextToProviderAssignment, contextId)
}

// UnregisterAndPop
// returns nil if there's not any more assignment
func (c *ContextProviderResolver) PopAndUnregister() (*string, *Provider) {
    // TODO: this feels a bit hacky to me, should look if there's a better
    // approach
    var selected string
    for k := range c.contextToProviderAssignment {
        selected = k
    }

    if selected == "" {
        return nil, nil
    }

    provider, ok := c.contextToProviderAssignment[selected];
    if !ok {
        panic("unexpected assignment not found after selection for popping")
    }

    c.Unregister(selected)

    return &selected, provider
}

func (c *ContextProviderResolver) UnregisterProviderEntries(p *Provider) {
	if p == nil {
		log.Fatalln("provider to have contexts deregistered should not be null")
	}

	// BUG: there's a race condition here, if i'm routing a new event
	// with not assignment at the smae time, the context could have been
	// created behind a region already passed by the loop
	// Should use a thread-safe data structure for this
	c.log("unregistering provider ", p.ID, "'s contexts")
	for k, v := range c.contextToProviderAssignment {
		if v.ID == p.ID {
			delete(c.contextToProviderAssignment, k)
		}
	}
}

func (c *ContextProviderResolver) TryRouteEvent(e *tooleventmodule.ToolEvent) error {
	if c.contextToProviderAssignment == nil {
		c.contextToProviderAssignment = map[string]*Provider{}
	}

	contextId := e.ContextId

	c.log("trying to route event from context: ", e.ContextId)
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

func (c *ContextProviderResolver) log(v ...any) {
	x := append([]any{"[ContextProviderResolver]"}, v...)

	log.Println(x...)
}
