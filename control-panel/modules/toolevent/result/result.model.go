package tooleventresult

type ToolEventResultStatus string

const (
	ToolEventResultStatusSuccess ToolEventResultStatus = "success"
	ToolEventResultStatusFailure ToolEventResultStatus = "failure"
)

type ToolEventResult struct {
	Status  ToolEventResultStatus `json:"status"`
	Message string                `json:"message"`
}

func (r *ToolEventResult) IsValid() bool {
	return r.Status == ToolEventResultStatusSuccess || r.Status == ToolEventResultStatusFailure
}
