package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func JsonPrint(v interface{}) {
	j, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(j))
}
