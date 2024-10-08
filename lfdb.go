package main

import (
	//"fmt"
	//"os"
	"net/http"
	"html/template"
)

type Person struct {
	FirstName 	string
	LastName	string
	DOB			string
	Occupation	string
}

type Film struct {
	FilmID		int
	Title		string
	Year		int
	Director	Person
	Genre		string
	Rating 		float32
	Plotpoints	[]string
	Series		bool
}

var film1 Film
var person1 Person
var tpl *template.Template

func main() {
	person1 = Person{
		FirstName: "Sergio",
		LastName: "Leone",
		DOB: "03.01.1929",
		Occupation: "Director",
	}

	film1 = Film{
		FilmID: 1,
		Title: "Once Upon a Time in The West",
		Year: 1968,
		Director: person1,
		Genre: "Western",
		Rating: 0.9,
		Plotpoints: []string{"harmonica", "widow", "former prostitute", "wild west", "spaghetti western"},
		Series: false,
	}

	tpl, _ = tpl.ParseGlob("webpage/*.html")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/film", filmHandler)

	http.ListenAndServe(":8080", nil)	
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
	//dat, _ := os.ReadFile("./test")
	//fmt.Fprintf(w, string(dat))
}

func filmHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "filmInfo.html", film1)
}
