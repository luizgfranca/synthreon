package tooleventmodule

import (
	tooleventdisplay "platformlab/controlpanel/modules/toolevent/display"
	tooleventinput "platformlab/controlpanel/modules/toolevent/input"
	tooleventresult "platformlab/controlpanel/modules/toolevent/result"
	"slices"
)

type EventType string

// BUG: when the final display command is sent from the server to the client there's no way
// for the client to know if the next reaction should be sent to the tool or not

// FIXME: in the future, when we want to reconnect, should we use the handshake ID to prove
// the provider is itself? Would this be unsafe?
//
// there also should be a way to retransmit events that happened in the disconnected period

// FIXME: when the client tries to reconnect, should we use the terminal_id and session to
// sinalize that he is hilself? If yes, there should be a new event type for the restart of
// the tool, and the events between the disconnected period should be transmitted back to it
//
// a question that stays is, is the user required to still have the exact state to reproduce
// the last tool screen before the disconnection or this should also be resent upon reconnection?

// FIXME: in both cases mentioned above, it the disconnected client or provider emits an
// event while is disconnected, it should try retransmitting it upon reconnection

// FIXME: between disconnections and reconnections there could be a interaction from the user upon
// a tool after an event has already been emitted by the provider's execution. A kind of causal
// reasoning should be implemented to avoid problems related to this.
//
// This can also happen when the provider emmits an event after a client emitted his but before
// having received it. Same kind of causal relationship and conflict resolution should be thought
//
// The current strictly interactive components don't have this problem, but in the future other
// components could have
const (
	// HANDSHAKE EVENT TYPES (between PROVIDER and SERVER)
	// a handshake event is the set of events that define the first contact between the
	// the provider and the server, in which the provider declares the project
	// it will provide tools for and the server assigns the provider an ID.

	// a handshake request is sent by the provider to the server identifying the
	// project and a handshake_id that will be used to correlate the interaction
	// the server then responds with a handshake/ack or a handshake/nack
	EventTypeHandshakeRequest EventType = "handshake/request"

	// success response in wich the server references the original handshake_id
	// and assigns the provider a provider_id that will be used to reference it back
	EventTypeHandshakeACK EventType = "handshake/ack"

	// erro response when it was not possible for any reason to register the provider
	// the reason will be sent in the "reason" parameter
	EventTypeHandshakeNACK EventType = "handshake/nack"

	// --
	// ANNOUNCEMENT EVENT TYPES (between PROVIDER and SERVER)
	// in an announcement event an already registered provider announces the tools
	// it intends to handle. Each announcement relates to one tool, if one wants to
	// handle more than one it should generate one event for each.

	// each announcement/handler should handle one tool handling request, it should have
	// the project and the tool it wants to handle, the tool should be registered in the
	// specified project.
	// The provider's provider_id should also be sent and the provider needs to also have
	// been registered in a handshake previously for the project.
	// The handshake_id is also needed for future-proofing correlations
	EventTypeAnnouncementHandler EventType = "announcement/handler"

	// the success response for the tool handling request.
	// it should have all the information contained in the corresponding announcement/handler event
	// and an assigned handler_id that will be used to assign future events for it
	EventTypeAnnouncementACK EventType = "announcement/ack"

	// the success response for the tool handling request.
	// it shuld have all the information contained in the corresponding announcement/handler event
	// and the reason the handling request could not be acknowleged
	EventTypeAnnouncementNACK EventType = "announcement/nack"

	// --
	// INTERACTION EVENT TYPES (between CLIENT and SERVER or SERVER and PROVIDER)
	// Represents an action from a user upon the tool it is using.
	// Each interaction is sent from the client to the server, and there is mapped to
	// the corresponding provider and handler that should deal with it, when this happens
	// the body of the data stays the smae, data from the requesting client and its execution
	// context is replaced by the provider and handler data

	// A user opens the tool, so the event should contain the project and the tool the user
	//
	// When from CLient to Server:
	// Just openned and, since they should already be logged in, their session_id, wich should
	// facilitate choosing to what communication channel to send the tool's reactions
	//
	// When from Server to Provider
	// Since the origin is from the client, this eventt is created from it.
	// The seession id is removed, since it should not need to be known by the server
	// and the provider_id and handler_id from the handler that should provide the tool
	// is added.
	EventTypeInteractionOpen EventType = "interaction/open"

	// FIX: does it really need the execution_id here?
	// FIX: to simplify things, shouldn't we remove the context_id?
	// An action from the user upon an already running tool
	//
	// When from Client to Server
	// It should indicate the project and tool, the context_id provided by the tool interactions
	// its terminal_id and the execution_id from the handler
	// it also should send the definition of the fields and values filled with the user interaction
	// according with the interface used in the step
	//
	// When from Server to Client
	// Event is made from the client event, but with both the context_id and the terminal_id removed
	// because the provider shouldn't need to know about then, and the filling of the provider_id and
	// handler_id to which the event should be sent, that should be the one that handled the user's
	// previous interactions
	EventTypeInteractionInput EventType = "interaction/input"

	// COMMAND EVENT TYPES (between CLIENT and SERVER or SERVER and PROVIDER)
	// Transmits a command
	// from the provider to the server
	// from the provider to the client (routed by the server)
	// from the server to the client
	//
	// the ones from the provider to the server should have the current handler's provider_id,
	// handler_id and execution_id,
	//
	// the ones from the server to the client should have the context_id of the run in the client
	// that should receive the command
	//
	// the ones from the provider to the client should behave the same way as if it was a command
	// from the provider to the server while passing between them, ant then as one between
	// server and client

	// A provider sinalizes that something should be shown to the user
	// This originates from the provider, and it passes through the server before being sent
	// to the client (with the required parameter transformations done)
	// This event type should have a "display" parameter that specifies which elements should
	// be shwon in the interface.
	EventTypeCommandDisplay EventType = "command/display"

	// A provider sinalizes to the server that the current execution for the tool has
	// been finished. The server then takes the appropriate steps in its state to finalize
	// the context and then sends a display event to the client to communicate this to the user.
	//
	// This event requires a "result" parameter that contains the status and details about why
	// this status was reached, as defined by the tool's implementation
	EventTypeCommandFinish EventType = "command/finish"
)

type ToolEvent struct {
	Type EventType `json:"type"`

	Project        string          `json:"project"`
	Tool           string          `json:"tool"`
	ToolProperties *ToolProperties `json:"tool_properties"`

	HandshakeId string `json:"announcement_id"`
	ProviderId  string `json:"provider_id"`
	HandlerId   string `json:"handler_id"`
	ExecutionId string `json:"execution_id"`
	TerminalId  string `json:"terminal_id"`

	SessionId string `json:"session_id"`
	ContextId string `json:"context_id"`

	Display *tooleventdisplay.DisplayDefinition `json:"display"`
	Input   *tooleventinput.InputDefinition     `json:"input"`
	Result  *tooleventresult.ToolEventResult    `json:"result"`

	Reason string `json:"reason"`
}

func (e *ToolEvent) IsHandshake() bool {
	return e.Type == EventTypeHandshakeRequest ||
		e.Type == EventTypeHandshakeACK ||
		e.Type == EventTypeHandshakeNACK
}

func (e *ToolEvent) IsAnnouncement() bool {
	return e.Type == EventTypeAnnouncementHandler ||
		e.Type == EventTypeAnnouncementACK ||
		e.Type == EventTypeAnnouncementNACK
}

func (e *ToolEvent) IsInteraction() bool {
	return e.Type == EventTypeInteractionOpen ||
		e.Type == EventTypeInteractionInput
}

func (e *ToolEvent) IsCommand() bool {
	return e.Type == EventTypeCommandDisplay ||
		e.Type == EventTypeCommandFinish
}

func (e *ToolEvent) IsValid() bool {
	if !slices.Contains([]EventType{
		EventTypeHandshakeRequest,
		EventTypeHandshakeACK,
		EventTypeHandshakeNACK,
		EventTypeAnnouncementHandler,
		EventTypeAnnouncementACK,
		EventTypeAnnouncementNACK,
		EventTypeInteractionOpen,
		EventTypeInteractionInput,
		EventTypeCommandDisplay,
		EventTypeCommandFinish,
	}, e.Type) {
		return false
	}

	return (e.IsHandshake() && e.passesHandshakeRules()) ||
		(e.IsAnnouncement() && e.passesAnnouncementRules()) ||
		(e.IsInteraction() && e.passesInteractionRules()) ||
		(e.IsCommand() && e.passesCommandRules())
}

func (e *ToolEvent) passesHandshakeRules() bool {
	if e.HandshakeId == "" || e.Project == "" {
		return false
	}

	if e.Type == EventTypeHandshakeACK && e.ProviderId == "" {
		return false
	}

	if e.Type == EventTypeHandshakeNACK && e.Reason == "" {
		return false
	}

	return true
}

func (e *ToolEvent) passesAnnouncementRules() bool {
	if e.HandshakeId == "" ||
		e.ProviderId == "" ||
		e.Project == "" ||
		e.Tool == "" {
		return false
	}

	if e.Type == EventTypeAnnouncementACK && e.HandlerId == "" {
		return false
	}

	if e.Type == EventTypeAnnouncementNACK && e.Reason == "" {
		return false
	}

	return true
}

func (e *ToolEvent) passesInteractionRules() bool {
	if e.Project == "" || e.Tool == "" {
		return false
	}

	if e.Type == EventTypeInteractionOpen {
		if e.SessionId == "" &&
			(e.ProviderId == "" || e.HandlerId == "" || e.ContextId == "") {
			return false
		}

		if e.SessionId != "" &&
			(e.ProviderId != "" || e.HandlerId != "" || e.ContextId != "") {
			return false
		}
	}

	if e.Type == EventTypeInteractionInput {
		if e.Input == nil || !e.Input.IsValid() {
			return false
		}

		if !((e.ExecutionId != "" && e.ContextId != "" && e.TerminalId != "") ||
			(e.ProviderId != "" && e.HandlerId != "" && e.ExecutionId != "")) {
			return false
		}
	}

	return true
}

func (e *ToolEvent) commandEventMatchesProviderToServerAttributes() bool {
	return (e.ProviderId != "" && e.HandlerId != "" && e.ExecutionId != "" && e.ContextId != "")
}

func (e *ToolEvent) commandEventMatchesServerToClientAttributes() bool {
	return (e.ProviderId == "" && e.HandlerId == "" && e.ExecutionId == "" && e.ContextId != "")
}

func (e *ToolEvent) passesCommandRules() bool {
	if e.Project == "" || e.Tool == "" {
		return false
	}

	if !(e.commandEventMatchesProviderToServerAttributes() ||
		e.commandEventMatchesServerToClientAttributes()) {
		return false
	}

	if e.Type == EventTypeCommandFinish &&
		(e.ProviderId == "" ||
			e.HandlerId == "" ||
			e.ExecutionId == "" ||
			e.Result == nil ||
			!e.Result.IsValid()) {
		return false
	}

	if e.Type == EventTypeCommandDisplay &&
		(e.Display == nil || !e.Display.IsValid()) {
		return false
	}

	return true
}
