package clientmodule

import (
	"log"
	"sync"
)

// TODO: this type of registry seems to be a patterrn,
// should consider refactoring into a reusable concept
type ContextTerminalResolver struct {
	registry sync.Map
}

func (r *ContextTerminalResolver) TryRegister(contextId string, term *Terminal) (ok bool) {
	return r.registry.CompareAndSwap(contextId, nil, term)
}

func (r *ContextTerminalResolver) Unregister(contextId string) {
	r.registry.Delete(contextId)
}

func (r *ContextTerminalResolver) Resolve(contextId string) *Terminal {
	v, ok := r.registry.Load(contextId)
	if !ok {
		return nil
	}

	term, ok := v.(*Terminal)
	if !ok {
		log.Fatalln("unexpected behavior: TerminalContextResolver's registry should only hold pointers to terminals")
	}

	return term
}
