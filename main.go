package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/Joker/jade"
	"github.com/namsral/flag"
)

var PORT = "3000"
var HOST string

type Page struct {
	OAuth2Link string
	Distance   float64
}

func Exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

type GeoUser struct {
	Email       string
	Name        string
	Ips         []string
	Geo         string
	Country     string
	Region_name string
	City        string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	hasToken := Exists("my_token")
	if hasToken {
		// render cool page
		buf, err := ioutil.ReadFile("views/index.jade")
		if err != nil {
			fmt.Printf("\nRead jade file error: %v", err)
			return
		}
		jadeTpl, err := jade.Parse("jade_tpl", string(buf))
		if err != nil {
			fmt.Printf("\nParse jade file error: %v", err)
			return
		}
		goTpl, err := template.New("html").Parse(jadeTpl)
		if err != nil {
			fmt.Printf("\nTemplate parse error: %v", err)
			return
		}
		if !Exists("data.json") {
			CreateDistanceData()
		}
		buf, err = ioutil.ReadFile("data.json")
		if err != nil {
			fmt.Printf("\nCould not read data.json")
			return
		}
		var data LocationData
		json.Unmarshal(buf, &data)
		goTpl.Execute(w, &Page{OAuth2Link: "", Distance: float64(int(data.Meters + data.ExtraMeters))})
	} else {
		// render login page
		link := GetOauth2Link()
		buf, err := ioutil.ReadFile("views/login.jade")
		if err != nil {
			fmt.Printf("\nRead jade file error: %v", err)
			return
		}
		jadeTpl, err := jade.Parse("jade_tpl", string(buf))
		if err != nil {
			fmt.Printf("\nParse jade file error: %v", err)
			return
		}
		goTpl, err := template.New("html").Parse(jadeTpl)
		if err != nil {
			fmt.Printf("\nTemplate parse error: %v", err)
			return
		}
		goTpl.Execute(w, &Page{OAuth2Link: link})
	}
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	values, _ := url.ParseQuery(r.URL.RawQuery)
	token := GetToken(values.Get("code"))
	if token == nil {
		return
	}
	mt := &myToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
	fmt.Printf("\naccess token: %s", mt.AccessToken)
	fmt.Printf("\ntoken type: %s", mt.TokenType)
	fmt.Printf("\nrefresh token: %s", mt.RefreshToken)
	fmt.Printf("\nexpires: %v", mt.Expiry)
	b, err := json.Marshal(mt)
	if err != nil {
		fmt.Println("Could not create JSON")
		return
	}
	err = ioutil.WriteFile("my_token", b, 0600)
	if err != nil {
		fmt.Println("Could not write JSON file")
		return
	}

	http.Redirect(w, r, "/", 308)
}

func main() {
	flag.StringVar(&HOST, "HOST", "localhost:"+PORT, "host for the oauth callback")
	flag.Parse()
	go FitnessRoutine()
	http.HandleFunc("/", viewHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/oauth2callback", callbackHandler)
	http.ListenAndServe(":"+PORT, nil)
}
