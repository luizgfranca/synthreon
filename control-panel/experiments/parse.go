package experiments

import (
	"encoding/json"
	"fmt"
)

type Test struct {
	Title string `json:"title"`
	Data  json.RawMessage
}

type Data struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

const (
	str = "{\"title\":\"test\",\"data\":[{\"name\":\"a\",\"value\":\"x\"},{\"name\":\"b\",\"value\":\"y\"}]}"
)

func Parse() {
	fmt.Println(str)

	var x Test
	json.Unmarshal([]byte(str), &x)

	fmt.Println(x)

	var items []Data
	json.Unmarshal(x.Data, &items)

	for i := range items {
		it := items[i]

		fmt.Println("item", it)
	}

}
