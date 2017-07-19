package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/fitness/v1"
)

const nanosPerMilli = 1e6

var (
	start        = time.Date(2017, time.April, 1, 0, 0, 0, 0, time.UTC)
	end          = time.Date(2017, time.July, 15, 0, 0, 0, 0, time.UTC)
	dataStreamId string
)

var config = &oauth2.Config{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://accounts.google.com/o/oauth2/token",
	},
	RedirectURL: "http://localhost:3000/oauth2callback",
	Scopes:      []string{"https://www.googleapis.com/auth/fitness.location.read"},
}

type secret struct {
	ClientID     string
	ClientSecret string
}

type myToken struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

func CreateConfig() {
	b, _ := ioutil.ReadFile("secret.json")
	var secret secret
	json.Unmarshal(b, &secret)
	config.ClientID = secret.ClientID
	config.ClientSecret = secret.ClientSecret
}

func GetOauth2Link() string {
	offline := oauth2.SetAuthURLParam("access_type", "offline")
	force := oauth2.SetAuthURLParam("approval_prompt", "force")
	return config.AuthCodeURL("", offline, force)
}

func GetToken(code string) *oauth2.Token {
	ctx := context.Background()
	token, err := config.Exchange(ctx, code)
	if err != nil {
		fmt.Printf("\ntoken error: %v\n", err)
		return nil
	}
	return token
}

func GetFitnessService() *fitness.Service {
	ctx := context.Background()
	b, _ := ioutil.ReadFile("my_token")
	var myToken myToken
	json.Unmarshal(b, &myToken)
	token := &oauth2.Token{
		AccessToken:  myToken.AccessToken,
		TokenType:    myToken.TokenType,
		RefreshToken: myToken.RefreshToken,
		Expiry:       myToken.Expiry,
	}
	client := config.Client(ctx, token)
	svc, err := fitness.New(client)
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
		return nil
	}
	if len(ds.DataSource) == 0 {
		fmt.Printf("\nYou have no user datasources to explore.", err)
		return nil
	}

	for _, v := range ds.DataSource {
		if strings.Contains(v.DataStreamId, "withings-distances") {
			dataStreamId = v.DataStreamId
		}
	}

	timeRange := fmt.Sprintf("%v-%v", start.UnixNano(), end.UnixNano())

	data, err := svc.Users.DataSources.Datasets.Get("me", dataStreamId, timeRange).Do()
	if err != nil {
		fmt.Printf("\nCan not read datasets: %v", err)
		return nil
	}

	return data
}

func LogDistanceData() {
	CreateConfig()
	t := time.NewTicker(time.Hour)
	waitForSecretFile := time.NewTicker(time.Minute)
	for {
		if !Exists("my_token") {
			fmt.Println("wait for secret to be created")
			<-waitForSecretFile.C
			continue
		}
		data := GetDistanceData()
		latestStartNano := start.UnixNano()
		kmLastDay := 0.0
		totalKm := 0.0
		// tm := time.Unix(0, latestStartNano)
		// day := tm.YearDay()

		for i, set := range data.Point {
			if i == 0 {
				latestStartNano = set.StartTimeNanos
				// fmt.Printf("\nsecond latestStartNano %d", latestStartNano)
				// tm := time.Unix(0, latestStartNano)
				// day = tm.YearDay()
			}
			if latestStartNano < set.StartTimeNanos {
				totalKm = totalKm + data.Point[i-1].Value[0].FpVal/1000
				// fmt.Printf("\n%d: %v km", day, data.Point[i-1].Value[0].FpVal/1000)
				// tm := time.Unix(0, latestStartNano)
				// day = tm.YearDay()
				latestStartNano = set.StartTimeNanos
			}
			if i == len(data.Point)-1 {
				kmLastDay = set.Value[0].FpVal / 1000
				// fmt.Printf("\n%d: %v km", day, kmLastDay)
			}
		}

		locationData := &LocationData{
			StartTimeNanos: latestStartNano,
			Meters:         totalKm,
			ExtraMeters:    kmLastDay,
		}

		SaveBulk(locationData)

		// fmt.Printf("\nTotal km: %v", totalKm+kmLastDay)
		<-t.C
	}
}
