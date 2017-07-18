package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type LocationData struct {
	StartTimeNanos int64
	Meters         float64
	ExtraMeters    float64
}

func SaveBulk(data *LocationData) {
	fmt.Printf("\nMeters: %v", data.Meters)
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("\nCould not create JSON: %v", err)
		return
	}
	err = ioutil.WriteFile("data.json", b, 0600)
	if err != nil {
		fmt.Println("Could not write JSON file")
		return
	}
}
