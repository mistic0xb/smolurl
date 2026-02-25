package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(msg string, v interface{}) {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	fmt.Printf("%s\n", fmt.Sprintf("-- %s --", msg))
	fmt.Println(string(json))
}
