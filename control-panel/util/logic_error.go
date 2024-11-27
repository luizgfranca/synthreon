package util

import "fmt"

type GenericLogicError struct {
	Message string
}

func (g *GenericLogicError) Error() string {
	return fmt.Sprintf("[GenericLogicError] %s", g.Message)
}
