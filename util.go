package main

import (
	"encoding/json"
	"fmt"
	"kindlenotes/kindle"
	"log"
)

// https://golangbyexample.com/print-struct-variables-golang/
func PrintClippingAsJson(clipping kindle.Section) {
	empJSON, err := json.Marshal(clipping)
	if err != nil {
		log.Fatalf(err.Error())
	}

	//MarshalIndent
	empJSON, err = json.MarshalIndent(clipping, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", string(empJSON))
}
