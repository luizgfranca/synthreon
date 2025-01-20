package contextmodule

type ExecutionContext struct {
	ID string

	ProviderId  string `json:"provider_id"`
	HandlerId   string `json:"handler_id"`
	ExecutionId string `json:"execution_id"`
	TerminalId  string `json:"terminal_id"`
}
