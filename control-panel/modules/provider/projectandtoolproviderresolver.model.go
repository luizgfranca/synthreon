package providermodule

import (
	"fmt"
	"log"
	"math/rand/v2"
	commonmodule "platformlab/controlpanel/modules/common"
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

func (p *ProjectAndToolProviderResolver) Resolve(project string, tool string) (*Provider, error) {
	key := key(project, tool)

	p.indexLock.Lock()
	v, ok := p.registry[key]
	p.indexLock.Unlock()
	if !ok {
		return nil, &commonmodule.GenericLogicError{Message: fmt.Sprintf("no tool found for <%s,%s>:", project, tool)}
	}

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
