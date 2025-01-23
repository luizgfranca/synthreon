package clientmodule

import (
	"log"
	"sync"
)

// TODO: this type of registry seems to be a patterrn,
// should consider refactoring into a reusable concept
type TerminalContextResolver struct {
	registry sync.Map
}

func (r *TerminalContextResolver) TryRegister(terminalId string, contextId string) (ok bool) {
	return r.registry.CompareAndSwap(terminalId, nil, contextId)
}

func (r *TerminalContextResolver) Unregister(terminalId string) {
	r.registry.Delete(terminalId)
}

func (r *TerminalContextResolver) Resolve(terminalId string) *string {
	v, ok := r.registry.Load(terminalId)
	if !ok {
		return nil
	}

	ctx, ok := v.(string)
	if !ok {
		log.Fatalln("unexpected behavior: TerminalContextResolver's registry should only hold strings")
	}

	return &ctx
}
