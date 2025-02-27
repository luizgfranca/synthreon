import { ToolEventDto } from 'platformlab-core';
import { DisplayRenderer, Field } from "@/component/displayRenderer";
import { useCallback } from 'react';
import { EmptyState } from './emptyState';

type CompositorProps = {
    contextId: string;
    terminalId: string;
    lastEvent: ToolEventDto;
    reset: () => void;
    sendEvent: (event: ToolEventDto) => void;
}

function empty(message: string) {
    return (
        <div className="flex pt-10 items-center">
            <EmptyState>{message}</EmptyState>
        </div>
    )
}

export default function Compositor(props: CompositorProps) {
    if (!props.contextId || !props.terminalId) {
        console.error('no contextId or no terminalId in context', props);
        return empty('Internal error: terminal not correctly setup');
    }

    const sendInteraction = useCallback((fields: Field[]) => {
        props.sendEvent({
            type: 'interaction/input',
            context_id: props.contextId,
            terminal_id: props.terminalId,
            input: { fields }
        })
    }, [props.sendEvent])

    return (
        <div className="text-zinc-100 h-screen pt-10">
            <div className="container mx-auto px-4">
                <DisplayRenderer
                    event={props.lastEvent}
                    onSumission={(fields) => sendInteraction(fields)}
                    resetCallback={props.reset}
                />
            </div>
        </div>
    )
}
