package main

import (
	"encoding/json"
	"english/cmd/domain/model"
	"fmt"
)

func main() {
	example := &model.Example{
		Example:     "test",
		Translation: "test",
	}

	m, _ := json.Marshal(example)

	fmt.Printf("abcd:%s", m)
}
