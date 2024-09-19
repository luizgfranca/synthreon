package util

import (
	"encoding/json"
	"fmt"
)

func Probe(item interface{}) {
	s, err := json.Marshal(item)
	if err != nil {
		fmt.Println("unable to probe struct")
		return
	}

	fmt.Println(string(s))
}
