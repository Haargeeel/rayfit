package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Joker/jade"
)

var PORT = "3000"

type Page struct {
	OAuth2Link string
	Distance   float64
}

func millisToTime(t int64) time.Time {
	return time.Unix(0, t*nanosPerMilli)
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
	// session, err := mgo.Dial("localhost")
	// if err != nil {
	// panic(err)
	// }
	// defer session.Close()

	// // Optional. Switch the session to a monotonic behavior.
	// session.SetMode(mgo.Monotonic, true)

	// c := session.DB("main").C("geo_user")
	// // err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
	// // &Person{"Cla", "+55 53 8402 8510"})
	// // if err != nil {
	// // log.Fatal(err)
	// // }

	// result := GeoUser{}
	// err = c.Find(bson.M{"first_name": "Ray"}).One(&result)
	// if err != nil {
	// log.Fatal(err)
	// }

	// for i, x := range result.Ips {
	// fmt.Printf("HALLO x: %s, i: %d\n", x, i)
	// }

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
	go LogDistanceData()
	http.HandleFunc("/", viewHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/oauth2callback", callbackHandler)
	http.ListenAndServe(":"+PORT, nil)
}
