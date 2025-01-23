package clientmodule

import (
	"log"
	"sync"
)

// TODO: this type of registry seems to be a patterrn,
// should consider refactoring into a reusable concept
type ContextClientResolver struct {
	registry sync.Map
}

func (r *ContextClientResolver) TryRegister(contextId string, term *Client) (ok bool) {
	return r.registry.CompareAndSwap(contextId, nil, term)
}

func (r *ContextClientResolver) Unregister(contextId string) {
	r.registry.Delete(contextId)
}

func (r *ContextClientResolver) Resolve(contextId string) *Client {
	v, ok := r.registry.Load(contextId)
	if !ok {
		return nil
	}

	term, ok := v.(*Client)
	if !ok {
		log.Fatalln("unexpected behavior: ContextClientResolver's registry should only hold pointers to clients")
	}

	return term
}
