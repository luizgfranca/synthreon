package tooleventinput

type InputDefinition struct {
	Fields []InputField `json:"fields"`
}

func (input *InputDefinition) IsValid() bool {
	for i := range input.Fields {
		if !input.Fields[i].IsValid() {
			return false
		}
	}

	return true
}
