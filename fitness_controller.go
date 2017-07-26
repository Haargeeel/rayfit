package main

import (
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/fitness/v1"
)

var (
	start        = time.Date(2017, time.April, 1, 0, 0, 0, 0, time.UTC)
	dataStreamId string
)

func GetFitnessService() *fitness.Service {
	svc, err := fitness.New(GetClient())
	if err != nil {
		fmt.Printf("\nCould not create fitness service: %v", err)
		return nil
	}
	return svc
}

func GetDistanceData() *fitness.Dataset {
	svc := GetFitnessService()
	ds, err := svc.Users.DataSources.List("me").Do()
	if err != nil {
		fmt.Printf("\nUnable to retrieve user's sessions: %v", err)
	}
	if len(ds.DataSource) == 0 {
		fmt.Printf("\nYou have no user datasources to explore.", err)
	}

	for _, v := range ds.DataSource {
		if strings.Contains(v.DataStreamId, "withings-distances") {
			dataStreamId = v.DataStreamId
		}
	}

	end := time.Now()

	timeRange := fmt.Sprintf("%v-%v", start.UnixNano(), end.UnixNano())

	data, err := svc.Users.DataSources.Datasets.Get("me", dataStreamId, timeRange).Do()
	if err != nil {
		fmt.Printf("\nCan not read datasets: %v", err)
	}

	return data
}

func CreateDistanceData() {
	data := GetDistanceData()
	latestStartNano := start.UnixNano()
	kmLastDay := 0.0
	totalKm := 0.0
	var days []*Day
	var lastDay *Day

	for i, set := range data.Point {
		if i == 0 {
			latestStartNano = set.StartTimeNanos
		}
		if latestStartNano < set.StartTimeNanos {
			totalKm = totalKm + data.Point[i-1].Value[0].FpVal/1000
			day := &Day{
				StartTime: latestStartNano,
				Meters:    data.Point[i-1].Value[0].FpVal,
			}
			days = append(days, day)
			latestStartNano = set.StartTimeNanos
		}
		if i == len(data.Point)-1 {
			lastDay = &Day{
				StartTime: latestStartNano,
				Meters:    set.Value[0].FpVal,
			}
			kmLastDay = set.Value[0].FpVal / 1000
		}
	}

	locationData := &LocationData{
		StartTimeNanos: latestStartNano,
		Meters:         totalKm,
		ExtraMeters:    kmLastDay,
	}

	SaveBulk(locationData)
	SaveDays(days, lastDay)
}

func FitnessRoutine() {
	t := time.NewTicker(time.Hour)
	waitForSecretFile := time.NewTicker(time.Minute)
	for {
		if !Exists("my_token") {
			fmt.Println("wait for secret to be created")
			<-waitForSecretFile.C
			continue
		}
		CreateDistanceData()
		<-t.C
	}
}
