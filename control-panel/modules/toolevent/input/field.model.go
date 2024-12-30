package tooleventinput

type InputField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (i *InputField) IsValid() bool {
	return i.Name != ""
}
