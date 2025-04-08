package providermodule

import (
	"fmt"
	"log"
	"math/rand/v2"
	commonmodule "synthreon/modules/common"
	"sync"
)

// TODO: the locking approach here seems inneficient to me,
// should try to improve it later
// FIXME: do deregistration
type ProjectAndToolProviderResolver struct {
	indexLock sync.Mutex
	registry  map[string]*slot
}

type slot struct {
	sync.Mutex
	list []*Provider
}

// TODO: maybe i should think of a more efficient index approach
func (p *ProjectAndToolProviderResolver) Register(
	project string,
	tool string,
	provider *Provider,
) {
	if p.registry == nil {
		p.registry = map[string]*slot{}
	}

	key := key(project, tool)
	log.Println("[ProjectAndToolProviderResolver] registering provider for key", key)

	p.indexLock.Lock()
	if p.registry[key] == nil {
		log.Println("[ProjectAndToolProviderResolver] creating slot")
		p.registry[key] = &slot{
			list: []*Provider{},
		}
	}
	p.indexLock.Unlock()

	item := p.registry[key]
	log.Println("[ProjectAndToolProviderResolver] appending provider to slot list, ", item)
	item.Lock()
	item.list = append(item.list, provider)
	item.Unlock()
}

func (p *ProjectAndToolProviderResolver) UnregisterProviderEntries(provider *Provider) {
	// TODO: think of a better structure to do this in a more efficient way
	for _, slot := range p.registry {
		for i, it := range slot.list {
			if it.ID == provider.ID {
				slot.list = commonmodule.RemoveFromUnorderedSlice(slot.list, i)

				// not exiting bere because the provider may be registered in
				// another item of the registry, so just exiting the slot loop here
				break
			}
		}
	}
}

func (p *ProjectAndToolProviderResolver) Resolve(project string, tool string) (*Provider, error) {
	key := key(project, tool)

	p.indexLock.Lock()
	v, ok := p.registry[key]
	p.indexLock.Unlock()
	if !ok {
		return nil, &commonmodule.GenericLogicError{Message: fmt.Sprintf("no tool found for <%s,%s>:", project, tool)}
	}

	if len(v.list) == 0 {
		// there was a handler registered for this combination here before, but it has since
		// been deregistered
		return nil, &commonmodule.GenericLogicError{Message: fmt.Sprintf("no handler currently registered for <%s,%s>:", project, tool)}
	}

	// FIXME: enable live loading on first run
	// TODO: should implement a better load balancing approach
	v.Lock()
	selected := rand.IntN(len(v.list))
	result := v.list[selected]
	v.Unlock()
	return result, nil
}

func key(project string, tool string) string {
	return project + ":" + tool
}
