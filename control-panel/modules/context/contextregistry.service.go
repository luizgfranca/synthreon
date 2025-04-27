package contextmodule

import (
	"log"
	"sync"
)

type ContextRegistry struct {
    registry map[string]*Context

    // TODO: should optimize this with a better concurrency approach
    registryLock sync.Mutex
}

func NewContextRegistry() *ContextRegistry {
    registry := make(map[string]*Context)
    c := ContextRegistry{
        registry: registry,
    }

    return &c
}

// Register: 
// if context already exists will replace current one,
// thread-safe,
// context is not nullable,
// presumed unfailable
func (c *ContextRegistry) Register(ctx *Context) {
    if (ctx == nil) {
        log.Fatalln("Register called for null context")
    }

    c.log("registering context ", ctx.ID, ctx)
    c.registryLock.Lock()
    c.registry[ctx.ID] = ctx
    c.registryLock.Unlock()
}

// Unregister:
// thread-safe,
// if contextId does not exist will do nothing
func (c *ContextRegistry) Unregister(ctxId string) {
    c.log("unregistering context ", ctxId)
    c.registryLock.Lock()
    delete(c.registry, ctxId)
    c.registryLock.Unlock()
}

// Get:
// thread-safe,
// if context is not found returns null
func (c *ContextRegistry) Get(ctxId string) *Context {
    c.registryLock.Lock()
    ctx, ok := c.registry[ctxId]
    c.registryLock.Unlock()

    if !ok {
        return nil
    }

    return ctx
}

func (c *ContextRegistry) log(v ...any) {
	x := append([]any{"[ContextRegistry] "}, v...)

	log.Println(x...)
}
