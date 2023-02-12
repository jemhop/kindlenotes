package main

import (
	"encoding/json"
	"fmt"
	"kindlenotes/kindle"
	"log"
)

type Model[T any] struct {
	Data []T
}

// https://golangbyexample.com/print-struct-variables-golang/
func PrintStructAsJson(clipping kindle.Section) {
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

// couldnt work out how to implement a function with any struct as an argument even though itd work for this usecase so this is necessary
func PrintBookAsJson(clipping kindle.Book) {
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
