package tooleventmodule

import (
	"encoding/json"
	"log"
	"platformlab/controlpanel/modules/commonmodule"
	"strings"
)

func parseEventV0String(versionstr *string, eventstr *string) (*ToolEvent, error) {
	log.Println("[EventParser] parsing v0 event:", *versionstr)
	var e ToolEvent

	err := json.Unmarshal([]byte(*eventstr), &e)
	if err != nil {
		log.Println("[EventParser] could not parse event string: ", err.Error())
		return nil, err
	}

	return &e, nil
}

func ParseEventString(input *string) (*ToolEvent, error) {
	if input == nil || *input == "" {
		return nil, &commonmodule.GenericLogicError{Message: "empty or null event string"}
	}

	parts := strings.SplitN(*input, "|", 2)
	if len(parts) != 2 {
		return nil, &commonmodule.GenericLogicError{Message: "version prefix could not be separated"}
	}

	version := parts[0]
	event := parts[1]

	if version == "v0.0" {
		return parseEventV0String(&version, &event)
	}

	return nil, &commonmodule.GenericLogicError{Message: "unknown event version"}
}
