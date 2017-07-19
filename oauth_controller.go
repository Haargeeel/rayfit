package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/fitness/v1"
)

var Config = &oauth2.Config{}

type myToken struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

func GetConfig() *oauth2.Config {
	if Config.ClientID != "" {
		return Config
	}

	buf, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		fmt.Printf("\nCannot read client_secret.json: %v", err)
	}

	Config, err = google.ConfigFromJSON(buf, fitness.FitnessLocationReadScope)
	if err != nil {
		fmt.Printf("\nCannot create OAuth2 Config: %v", err)
	}

	Config.RedirectURL = "http://" + HOST + "/oauth2callback"

	return Config
}

func GetOauth2Link() string {
	return GetConfig().AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func GetClient() *http.Client {
	ctx := context.Background()

	buf, _ := ioutil.ReadFile("my_token")
	var myToken myToken
	json.Unmarshal(buf, &myToken)
	token := &oauth2.Token{
		AccessToken:  myToken.AccessToken,
		TokenType:    myToken.TokenType,
		RefreshToken: myToken.RefreshToken,
		Expiry:       myToken.Expiry,
	}
	return GetConfig().Client(ctx, token)
}

func GetToken(code string) *oauth2.Token {
	ctx := context.Background()
	token, err := GetConfig().Exchange(ctx, code)
	if err != nil {
		fmt.Printf("\ntoken error: %v\n", err)
	}
	return token
}
