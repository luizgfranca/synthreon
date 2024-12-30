package tooleventmodule

import (
	"fmt"
	tooleventdisplay "platformlab/controlpanel/modules/toolevent/display"
	tooleventresult "platformlab/controlpanel/modules/toolevent/result"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TODO: create tests for the validations of the data definitions (input, display, ...)
func TestToolEventRules(t *testing.T) {
	invalidDisplayExample := tooleventdisplay.DisplayDefinition{
		Type: tooleventdisplay.DisplayTypeInformation,
	}

	validDisplayExample := tooleventdisplay.DisplayDefinition{
		Type: tooleventdisplay.DisplayTypeInformation,
		Information: &tooleventdisplay.InformationDisplay{
			Type:    tooleventdisplay.InformationDisplayTypeSuccess,
			Message: "msg",
		},
	}

	invalidResultExample := tooleventresult.ToolEventResult{
		Status:  "fskjhdfkshjd",
		Message: "message",
	}

	validReusltExample := tooleventresult.ToolEventResult{
		Status:  tooleventresult.ToolEventResultStatusSuccess,
		Message: "message",
	}

	cases := []struct {
		input    ToolEvent
		expected bool
	}{
		{input: ToolEvent{}, expected: false},

		// PROVIDER NEGOTIATIONS
		{input: ToolEvent{Type: "handshake/request", Project: "sdfsdf", HandshakeId: "adfkhaksdhjf"}, expected: true},
		{input: ToolEvent{Type: "handshake/request", Project: "sdfsdf"}, expected: false},
		{input: ToolEvent{Type: "handshake/request"}, expected: false},
		{input: ToolEvent{Type: "handshake/ack"}, expected: false},
		{input: ToolEvent{Type: "handshake/ack", Project: "a"}, expected: false},
		{input: ToolEvent{Type: "handshake/ack", Project: "a", HandshakeId: "handsh"}, expected: false},
		{input: ToolEvent{Type: "handshake/ack", Project: "a", ProviderId: "prov"}, expected: false},
		{input: ToolEvent{Type: "handshake/ack", Project: "a", HandshakeId: "handsh", ProviderId: "prov"}, expected: true},
		{input: ToolEvent{Type: "handshake/nack", Project: "a", HandshakeId: "handsh", ProviderId: "prov"}, expected: false},
		{input: ToolEvent{Type: "handshake/nack", Project: "a", HandshakeId: "handsh"}, expected: false},
		{input: ToolEvent{Type: "handshake/nack", Project: "a", HandshakeId: "handsh", Reason: "reason"}, expected: true},

		{input: ToolEvent{Type: "announcement/handler", Project: "a", HandshakeId: "handsh"}, expected: false},
		{input: ToolEvent{Type: "announcement/handler", Project: "a", ProviderId: "prov"}, expected: false},
		{input: ToolEvent{Type: "announcement/handler", Project: "a", HandshakeId: "handsh", ProviderId: "prov"}, expected: false},
		{input: ToolEvent{Type: "announcement/handler", Project: "a", Tool: "t", HandshakeId: "handsh", ProviderId: "prov"}, expected: true},
		{input: ToolEvent{Type: "announcement/ack", Project: "a", HandshakeId: "handsh", ProviderId: "prov"}, expected: false},
		{input: ToolEvent{Type: "announcement/ack", Project: "a", HandshakeId: "handsh", ProviderId: "prov", HandlerId: "handl"}, expected: false},
		{input: ToolEvent{Type: "announcement/ack", Project: "a", Tool: "t", HandshakeId: "handsh", ProviderId: "prov", HandlerId: "handl"}, expected: true},
		{input: ToolEvent{Type: "announcement/nack", Project: "a", HandshakeId: "handsh", ProviderId: "prov", HandlerId: "handl"}, expected: false},
		{input: ToolEvent{Type: "announcement/nack", Project: "a", HandshakeId: "handsh", ProviderId: "prov"}, expected: false},
		{input: ToolEvent{Type: "announcement/nack", Project: "a", HandshakeId: "handsh", ProviderId: "prov", Reason: "reasn"}, expected: false},
		{input: ToolEvent{Type: "announcement/nack", Project: "a", Tool: "t", HandshakeId: "handsh", ProviderId: "prov", Reason: "reasn"}, expected: true},

		// OPEN INTERACITON
		// client to server
		{input: ToolEvent{Type: "interaction/open"}, expected: false},
		{input: ToolEvent{Type: "interaction/open", Project: "a"}, expected: false},
		{input: ToolEvent{Type: "interaction/open", Project: "a", Tool: "t"}, expected: false},
		{input: ToolEvent{Type: "interaction/open", Project: "a", Tool: "t", SessionId: "sess"}, expected: true},

		// server to provider
		{input: ToolEvent{Type: "interaction/open", Project: "a", Tool: "t", ProviderId: "pid"}, expected: false},
		{input: ToolEvent{Type: "interaction/open", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl"}, expected: false},
		{input: ToolEvent{Type: "interaction/open", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ContextId: "ctx"}, expected: true},

		// shouldnt mix them up
		{input: ToolEvent{Type: "interaction/open", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ContextId: "ctx", SessionId: "sess"}, expected: false},

		// INPUT INTERACTION
		{input: ToolEvent{Type: "interaction/input"}, expected: false},
		{input: ToolEvent{Type: "interaction/open", Project: "a"}, expected: false},
		{input: ToolEvent{Type: "interaction/open", Project: "a", Tool: "t"}, expected: false},

		// client to server
		{input: ToolEvent{Type: "interaction/open", Project: "a", Tool: "t", ExecutionId: "exec"}, expected: false},
		{input: ToolEvent{Type: "interaction/open", Project: "a", Tool: "t", ExecutionId: "exec", ContextId: "ctx"}, expected: false},
		{input: ToolEvent{Type: "interaction/open", Project: "a", Tool: "t", ExecutionId: "exec", ContextId: "ctx", TerminalId: "term"}, expected: false},

		// DISPLAY COMMAND
		{input: ToolEvent{Type: "command/display"}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a"}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t"}, expected: false},

		// provider to server
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", ProviderId: "pid"}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl"}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ContextId: "ctx"}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ContextId: "ctx", ExecutionId: "exec"}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ContextId: "ctx", ExecutionId: "exec", Display: &invalidDisplayExample}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ContextId: "ctx", ExecutionId: "exec", Display: &validDisplayExample}, expected: true},
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ContextId: "ctx", Display: &validDisplayExample}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ExecutionId: "exec", Display: &validDisplayExample}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", HandlerId: "handl", ContextId: "ctx", ExecutionId: "exec", Display: &validDisplayExample}, expected: false},

		// server to client
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", ContextId: "ctx"}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a", ContextId: "ctx"}, expected: false},
		{input: ToolEvent{Type: "command/display", Tool: "t", ContextId: "ctx"}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", ContextId: "ctx", Display: &invalidDisplayExample}, expected: false},
		{input: ToolEvent{Type: "command/display", Project: "a", Tool: "t", ContextId: "ctx", Display: &validDisplayExample}, expected: true},

		// FINISH COMMAND
		{input: ToolEvent{Type: "command/finish"}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a"}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", Tool: "t"}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", Tool: "t", ProviderId: "pid"}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl"}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ExecutionId: "exec"}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ExecutionId: "exec", Result: &invalidResultExample}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ExecutionId: "exec", Result: &validReusltExample}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ExecutionId: "exec", ContextId: "ctx", Result: &validReusltExample}, expected: true},
		{input: ToolEvent{Type: "command/finish", Tool: "t", ProviderId: "pid", HandlerId: "handl", ExecutionId: "exec", Result: &validReusltExample}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", ProviderId: "pid", HandlerId: "handl", ExecutionId: "exec", Result: &validReusltExample}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", ExecutionId: "exec", Result: &validReusltExample}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", Tool: "t", HandlerId: "handl", ExecutionId: "exec", Result: &validReusltExample}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", Tool: "t", ProviderId: "pid", ExecutionId: "exec", Result: &validReusltExample}, expected: false},
		{input: ToolEvent{Type: "command/finish", Project: "a", Tool: "t", ProviderId: "pid", HandlerId: "handl", Result: &validReusltExample}, expected: false},
	}

	fmt.Println("Test ToolEvent structure validations")
	for i := range cases {
		expectation := "should succeed"
		if !cases[i].expected {
			expectation = "should not succeed"
		}

		name := fmt.Sprintln(cases[i].input, expectation)
		t.Run(name, func(t *testing.T) {
			if cases[i].input.IsValid() != cases[i].expected {
				t.Error(cases[i].input, expectation)
			}
		})
	}
}

func TestToolEventParsing(t *testing.T) {
	fmt.Println("Test ToolEvent structure validations")

	cases := []struct {
		input          string
		expectedValid  bool
		expectedError  bool
		expectedResult *ToolEvent
	}{
		{
			input: `v0.0|{
				"type": "interaction/open",
				"project": "x",
				"tool": "y",
				"session_id": "shdfkajhdflaksjdfhalksdjh"
			}`,
			expectedValid: true,
			expectedResult: &ToolEvent{
				Type:      EventTypeInteractionOpen,
				Project:   "x",
				Tool:      "y",
				SessionId: "shdfkajhdflaksjdfhalksdjh",
			},
		}, {
			input: `v0.0|{
				"type": "command/display",
				"project": "x",
				"tool": "y",
				"context_id": "asdfgasdfhaskdfhj",
				"display": {
					"type": "information",
					"information": {
						"type": "success",
						"message": "Hello World"
					}
				}
			}`,
			expectedValid: true,
			expectedResult: &ToolEvent{
				Type:      EventTypeCommandDisplay,
				Project:   "x",
				Tool:      "y",
				ContextId: "asdfgasdfhaskdfhj",
				Display: &tooleventdisplay.DisplayDefinition{
					Type: tooleventdisplay.DisplayTypeInformation,
					Information: &tooleventdisplay.InformationDisplay{
						Type:    tooleventdisplay.InformationDisplayTypeSuccess,
						Message: "Hello World",
					},
				},
			},
		},
		{
			input: `v0.0|{
				"type": "command/display",
				"project": "x",
				"tool": "y",
				"context_id": "asdfgasdfhaskdfhj",
				"display": {
					"type": "information",
					"information": {
						"type": "error",
						"message": "Hello World"
					}
				}
			}`,
			expectedValid: false,
			expectedResult: &ToolEvent{
				Type:      EventTypeCommandDisplay,
				Project:   "x",
				Tool:      "y",
				ContextId: "asdfgasdfhaskdfhj",
				Display: &tooleventdisplay.DisplayDefinition{
					Type: tooleventdisplay.DisplayTypeInformation,
					Information: &tooleventdisplay.InformationDisplay{
						Type:    "error",
						Message: "Hello World",
					},
				},
			},
		}, {
			input: `|{
				"type": "interaction/open",
				"project": "x",
				"tool": "y",
				"session_id": "shdfkajhdflaksjdfhalksdjh"
			}`,
			expectedError: true,
		},
		{
			input: `{
				"type": "interaction/open",
				"project": "x",
				"tool": "y",
				"session_id": "shdfkajhdflaksjdfhalksdjh"
			}`,
			expectedError: true,
		},
	}

	for i := range cases {
		name := fmt.Sprintln(cases[i].input, cases[i].expectedValid, cases[i].expectedResult)
		t.Run(name, func(t *testing.T) {
			data, err := ParseEventString(&cases[i].input)

			if err != nil {
				if !cases[i].expectedError {
					t.Error(err)
				}

				return
			}

			if data.IsValid() != cases[i].expectedValid {
				t.Error("unexpected validation result expectation")
				return
			}

			if !cmp.Equal(*data, *cases[i].expectedResult) {
				t.Error(
					"parsed event result does not match\n Received:", *data,
					"\n Expected:", *cases[i].expectedResult,
					"\n", cmp.Diff(*data, *cases[i].expectedResult),
				)
				return
			}
		})
	}

}
