package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/Joker/jade"
)

var PORT = "80"

type Page struct {
	Title string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadFile("views/index.jade")
	if err != nil {
		fmt.Printf("\nReadfile error: %v", err)
		return
	}
	jadeTpl, err := jade.Parse("jade_tpl", string(buf))
	if err != nil {
		fmt.Printf("\nParse error: %v", err)
		return
	}
	goTpl, err := template.New("html").Parse(jadeTpl)
	if err != nil {
		fmt.Printf("\nTemplate parse error: %v", err)
		return
	}
	p := &Page{Title: "test"}
	goTpl.Execute(w, p)
}

func main() {
	fmt.Printf("Server listening on port: %s", PORT)
	http.HandleFunc("/", viewHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":"+PORT, nil)
}
