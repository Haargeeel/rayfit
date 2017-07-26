package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type LocationData struct {
	StartTimeNanos int64
	Meters         float64
	ExtraMeters    float64
}

type Day struct {
	StartTime int64   `json:"start_time"`
	Meters    float64 `json:"meters"`
}

type DayData struct {
	LastDay *Day   `json:"last_day"`
	Days    []*Day `json:"days"`
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

func SaveDays(days []*Day, lastDay *Day) {
	data := &DayData{
		LastDay: lastDay,
		Days:    days,
	}
	for _, v := range days {
		t := time.Unix(0, v.StartTime)
		_, month, date := t.Date()
		fmt.Printf("\n%v.%v.: %v meters", date, month, v.Meters)
	}
	t := time.Unix(0, lastDay.StartTime)
	_, month, date := t.Date()
	fmt.Printf("\n%v.%v.: %v meters", date, month, lastDay.Meters)
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("\nCould not create JSON: %v", err)
		return
	}
	err = ioutil.WriteFile("days.json", b, 0600)
	if err != nil {
		fmt.Println("Could not write JSON file")
		return
	}
}
