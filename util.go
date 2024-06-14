package main

import (
	"encoding/json"
	"log"
)

// printJSON 打印 JSON 格式的数据
func printJSON(data interface{}) {
	if *verbose{
		jsonBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			return
		}
	log.Printf("\nJSON:\n%s\n", string(jsonBytes))
	}
}

// print 打印数据
func print(data interface{}){
	if *verbose{
	log.Println(data)
	}
}
