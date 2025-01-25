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

func (r *ContextTerminalResolver) TryRegister(ctxid string, term *Terminal) (success bool) {
	_, ok := r.registry.Load(ctxid)
	if ok {
		return false
	}

	r.registry.Store(ctxid, term)
	return true
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
