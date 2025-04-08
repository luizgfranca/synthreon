package providermodule

import (
	"log"
	commonmodule "synthreon/modules/common"
	tooleventmodule "synthreon/modules/toolevent"
	"sync"
	"time"
)

type dlqSlot struct {
	event     *tooleventmodule.ToolEvent
	createdAt time.Time
}

// In the event dead letter queue will be saved events that
// could not be sent to their respective providers for one of
// two reasons
//   - there's no tool/project's provider currently registered
//     to handle the event
//   - the provider for the context created chould not be found
//     TODO: the behavior of this one will be implemented later
//     when a behavior for provider reconnection is implemented
//
// When an event fails it should be saved here, and when there's
// a behavior that inteacts with a project/tool registration
// or with a context, the DQL should be rerun for this and will
// check if it interacted with the context or project/tool of
// any of the events here, and if so, it should try to send them
// again.

// FIXME: implement sweeping
type EventDLQ struct {
	timeoutSeconds int

	ptmtx sync.Mutex
	// for any operation using this one should hold "ptmtx"
	projectAndToolSlots map[string]*[]dlqSlot
}

// newEventDLQ:
// return is never null
func newEventDLQ(timeoutSeconds int) *EventDLQ {
	q := EventDLQ{
		timeoutSeconds:      timeoutSeconds,
		projectAndToolSlots: make(map[string]*[]dlqSlot),
	}

	return &q
}

// sweepService:
// internal use,
// perform the sweep from sweepService()
func (q *EventDLQ) sweep() {
	for _, v := range q.projectAndToolSlots {
		for i, it := range *v {
			if time.Since(it.createdAt) >= time.Duration(q.timeoutSeconds)*time.Second {
				q.log("cleaning up expired event", it.event)
				*v = commonmodule.RemoveFromOrderedSlice(*v, i)
			}
		}
	}
}

// sweepService:
// should be called as a goroutine to run as a daemon in the background
func (q *EventDLQ) sweepService() {
	for {
		time.Sleep(100 * time.Millisecond)

		q.ptmtx.Lock()
		q.sweep()
		q.ptmtx.Unlock()
	}
}

// popSlot:
// internal use,
// should maintain the ordering of events,
// queue is not nullable
func (q *EventDLQ) popSlot(queue *[]dlqSlot) *tooleventmodule.ToolEvent {
	if queue == nil {
		log.Fatalln("queue to pop should never be null")
	}

	q.log("pop slot called; queue: ", queue)

	if len(*queue) == 0 {
		return nil
	}

	// TODO: consider turning queue behaviors into common utilities
	e := (*queue)[0].event
	*queue = (*queue)[1:]

	return e
}

// pushSlot:
// internal use,
// queue is not nullable,
// slot is not nullable
func (q *EventDLQ) pushSlot(queue *[]dlqSlot, slot *dlqSlot) {
	if queue == nil || slot == nil {
		log.Fatalln("queue or slot should never be null in push operation")
	}

	*queue = append(*queue, *slot)
}

// register:
// event is not nullable
func (q *EventDLQ) register(e *tooleventmodule.ToolEvent) {
	if e == nil {
		log.Fatalln("event should not be null while registering in DLQ")
	}

	s := dlqSlot{
		event:     e,
		createdAt: time.Now(),
	}

	k := q.ptkey(e.Project, e.Tool)

	q.log("registering dead letter queue entry: key:", k, ", value:", s)
	q.ptmtx.Lock()
	defer q.ptmtx.Unlock()

	queue, ok := q.projectAndToolSlots[k]
	if !ok {
		queue = new([]dlqSlot)
		q.projectAndToolSlots[k] = queue
	}

	q.pushSlot(queue, &s)
	q.log("queue after push:", queue)
}

// popFromProjectAndTool:
// return event is nullable
func (q *EventDLQ) popFromProjectAndTool(
	projectAcronym string, toolAcronym string,
) *tooleventmodule.ToolEvent {
	k := q.ptkey(projectAcronym, toolAcronym)
	q.log("popping from dead letter queue", k)

	q.ptmtx.Lock()
	defer q.ptmtx.Unlock()

	queue, ok := q.projectAndToolSlots[k]
	if !ok {
		return nil
	}

	return q.popSlot(queue)
}

func (q *EventDLQ) ptkey(projectAcronym string, toolAcronym string) string {
	return projectAcronym + ":" + toolAcronym
}

func (q *EventDLQ) log(v ...any) {
	x := append([]any{"[EventDLQ]"}, v...)

	log.Println(x...)
}
