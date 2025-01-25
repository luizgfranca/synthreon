import { DisplayDefinition } from "./display/display.dto";
import { InputDefinition } from "./input/input.dto";
import { ToolEventResult } from "./result/result.dto";

export const EventTypeValue = {
    HandshakeRequest: "handshake/request",
    HandshakeACK: "handshake/ack",
    HandshakeNACK: "handshake/nack",

    AnnouncementHandler: "announcement/handler",
    AnnouncementACK: "announcement/ack",
    AnnouncementNACK: "announcement/nack",

    InteractionOpen: "interaction/open",
    InteractionInput: "interaction/input",

    CommandDisplay: "command/display",
    CommandFinish: "command/finish",
} as const

type EventType = typeof EventTypeValue[keyof typeof EventTypeValue]

export type ToolEventDto = {
    type: EventType;
    project: string;
    tool: string;

    announcement_id?: string;
    provider_id?: string;
    handler_id?: string;
    execution_id?: string;
    terminal_id?: string;
    session_id?: string;
    context_id?: string;

    display?: DisplayDefinition;
    input?: InputDefinition;
    result?: ToolEventResult;

    reason?: string;
}