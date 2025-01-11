package commonmodule

import "fmt"

// TODO: In most cases this is used, a specialized structure should be created instead
type GenericLogicError struct {
	Message string
}

func (g *GenericLogicError) Error() string {
	return fmt.Sprintf("[GenericLogicError] %s", g.Message)
}
