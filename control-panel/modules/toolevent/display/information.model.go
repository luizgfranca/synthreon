package tooleventdisplay

import (
	"slices"
)

type InformationDisplayType string

const (
	InformationDisplayTypeSuccess InformationDisplayType = "success"
	InformationDisplayTypeFailure InformationDisplayType = "failure"
)

type InformationDisplay struct {
	Type    InformationDisplayType
	Message string
}

func (i *InformationDisplay) IsValid() bool {
	return slices.Contains(
		[]InformationDisplayType{InformationDisplayTypeFailure, InformationDisplayTypeSuccess},
		i.Type,
	) && i.Message != ""
}
