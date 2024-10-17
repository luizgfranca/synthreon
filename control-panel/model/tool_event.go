package model

type EventClass string

const (
	EventClassOperation    EventClass = "operation"
	EventClassInteraction  EventClass = "interaction"
	EventClassAnnouncement EventClass = "announcement"
)

type EventType string

const (
	// only used when eventClass=operation
	EventTypeDisplay EventType = "display"
	EventTypeInput   EventType = "input"

	// only used when eventClass=interaction
	EventTypeOpen EventType = "open"

	// only used when eventClass=announcement
	EventTypeProvider EventType = "provider"
	EventTypeAck      EventType = "ack"
)

type DisplayDefniitionType string

const (
	DisplayDefniitionTypeResult DisplayDefniitionType = "result"
	DisplayDefniitionTypeView   DisplayDefniitionType = "view"
	DisplayDefniitionTypePrompt DisplayDefniitionType = "prompt"
)

type DisplayElement struct {
	Type        string `json:"type"`
	Label       string `json:"label"`
	Text        string `json:"text"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type DisplayResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// TODO: thinking about this better, maybe i dont need the result,
//
//	should try to emulate it with specific elements
//	think about refactoring this
type DisplayDefniition struct {
	Type DisplayDefniitionType `json:"type"`

	// if type=(view or display) Elements should be present
	// if type=result Result should be present
	Elements *[]DisplayElement `json:"elements"`
	Result   *DisplayResult    `json:"result"`
}

type UserInput struct {
	// should always be present if defined
	Fields *[]interface{} `json:"fields"`
}

type ToolEvent struct {
	Class   EventClass `json:"class"`
	Type    EventType  `json:"type"`
	Project string     `json:"project"`
	Tool    string     `json:"tool"`

	// should be present if type=operation and type=display
	Display *DisplayDefniition `json:"display"`

	// should be present if type=interaction and type=input
	Input *UserInput `json:"input"`
}

func (e *ToolEvent) IsValid() bool {

	// validation of required fields

	if e.Class == "" || e.Type == "" || e.Project == "" || e.Tool == "" {
		return false
	}

	// validation of all possible event classes

	if e.Class != EventClassInteraction && e.Class != EventClassOperation && e.Class != EventClassAnnouncement {
		return false
	}

	// validation of all possible types

	if e.Type != EventTypeInput &&
		e.Type != EventTypeDisplay &&
		e.Type != EventTypeOpen &&
		e.Type != EventTypeProvider &&
		e.Type != EventTypeAck {
		return false
	}

	// validation of class/type possible combinations

	if e.Class == EventClassOperation &&
		(e.Type != EventTypeDisplay && e.Type != EventTypeInput) {
		return false
	}

	if e.Class == EventClassInteraction && e.Type != EventTypeOpen {
		return false
	}

	if e.Class == EventClassAnnouncement &&
		(e.Type != EventTypeProvider && e.Type != EventTypeAck) {
		return false
	}

	// validation of field configuration when specific class/type fields are needed

	if e.Class == EventClassOperation && e.Type == EventTypeDisplay {
		if e.Display == nil {
			return false
		}

		if (e.Display.Type == DisplayDefniitionTypePrompt ||
			e.Display.Type == DisplayDefniitionTypeView) &&
			e.Display.Elements == nil {
			return false
		}

		if e.Display.Type == DisplayDefniitionTypeResult && e.Display.Result == nil {
			return false
		}
	}

	if e.Class == EventClassInteraction &&
		e.Type == EventTypeInput &&
		(e.Input == nil || e.Input.Fields == nil) {
		return false
	}

	return true
}

// EXPECTATION: the event sent to this function should be an announcement
func NewAnnouncementAckEvent(announcement *ToolEvent) *ToolEvent {
	if announcement.Class != EventClassAnnouncement {
		panic("expected Announcement event in NewAnnouncementAckEvent() but received another Class: " + announcement.Class)
	}

	return &ToolEvent{
		Class:   EventClassAnnouncement,
		Type:    EventTypeAck,
		Project: announcement.Project,
		Tool:    announcement.Tool,
	}
}
