package main

import (
	"fmt"
	"net/http"
	"html/template"
	"io/ioutil"
	"strings"
	"encoding/json"
)

/*
TODO:
... proper functionality, database integration
 -> write settings.html to accept input apikey and filmdir (forms)
*/

/*
type FilmData struct {
	Title           string `json:"Title"` 		
	Year            string `json:"Year"`		
	Rated           string `json:"Rated"`		
	Released        string `json:"Released"`	
	Runtime			string `json:"Runtime"`		
	Genre           string `json:"Genre"`		
	Director        string `json:"Director"`	
	Writer          string `json:"Writer"`		
	Actors          string `json:"Actors"`		
	Plot            string `json:"Plot"`		
	Language        string `json:"Language"`	
	Country         string `json:"Country"`		
	Awards          string `json:"Awards"`		
	Poster          string `json:"Poster"`		
	Ratings  		string `json:"Ratings"`	
	Metascore		string `json:"Metascore`	
	ImdbRating		string `json:"imdbRating`	
	ImdbVotes		string `json:"imdbVotes"`	
	ImdbID  		string `json:"imdbID"`		
	Type			string `json:"Type"`		
	DVD				string `json:"DVD"`			
	BoxOffice		string `json:"BoxOffice"`	
	Production		string `json:"Production"`	
	Website			string `json:"Website"`		
	Response        string `json:"Response"`	
}

{"Title":"Once Upon a Time in the West",
"Year":"1968",
"Rated":"PG-13",
"Released":"04 Jul 1969",
"Runtime":"166 min",
"Genre":"Drama, Western",
"Director":"Sergio Leone",
"Writer":"Sergio Donati, Sergio Leone, Dario Argento",
"Actors":"Henry Fonda, Charles Bronson, Claudia Cardinale",
"Plot":"A mysterious stranger with a harmonica joins forces with a notorious desperado to protect a beautiful widow from a ruthless assassin working for the railroad.",
"Language":"Italian, English, Spanish",
"Country":"Italy, United States",
"Awards":"6 wins & 5 nominations",
"Poster":"https://m.media-amazon.com/images/M/MV5BZjYyNGY1MDEtN2I1MC00MGVhLTljZTYtODQ1NzQ0ODc2NzZlXkEyXkFqcGc@._V1_SX300.jpg",
"Ratings":[{"Source":"Internet Movie Database","Value":"8.5/10"},{"Source":"Rotten Tomatoes","Value":"96%"},{"Source":"Metacritic","Value":"82/100"}],
"Metascore":"82",
"imdbRating":"8.5",
"imdbVotes":"354,024",
"imdbID":"tt0064116",
"Type":"movie",
"DVD":"N/A",
"BoxOffice":"$5,321,508",
"Production":"N/A",
"Website":"N/A",
"Response":"True"}
*/

type Film struct {
	ImdbID		string
	Title		string
	Year		string
	Runtime		string
	Director	string
	Writer		string
	Genre		string
	Country		string
	Rating		string
	Plot		string
	Type		string
}

type SearchResult struct {
	Search []SearchMovie `json:"Search"`
}

type SearchMovie struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	IMDbID string `json:"imdbID"`
	Type   string `json:"Type"`
}
type SearchError struct {
	Response	string	`json:"Response"`
	Error		string 	`json:"Error"`
}

type FilmData struct {
	Title       string `json:"Title"` 		
	Year        string `json:"Year"`		
	Runtime		string `json:"Runtime"`		
	Genre       string `json:"Genre"`		
	Director    string `json:"Director"`	
	Writer		string `json:"Writer"`
	Actors      string `json:"Actors"`			
	Plot        string `json:"Plot"`		
	Country     string `json:"Country"`		
	Poster      string `json:"Poster"`		
	ImdbID      string `json:"imdbID"`	
	ImdbRating	string `json:"imdbRating`	
	Type		string `json:"Type"`		
	Response    string `json:"Response"`	
}

var film Film
var tpl *template.Template
var apikey string
var filmdir string

func main() {
	tpl, _ = tpl.ParseGlob("webpage/*.html")
	apikey = ""
	filmdir = "K:/Films"
	films := directoryReader()

	wf := len(films)-46

	fmt.Println("Wanted film:", films[wf])
	a := getFilmID(films[wf]) 
	getFilmData(a)
	//fmt.Println(b)

	
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
	tpl.ExecuteTemplate(w, "filmInfo.html", film)
}

func getFilmID(searchQuery string) (imdbID string) {// searches for a film and returns imdbID
	var searchResults SearchResult
	
	useCase := "s"
	title := searchQuery[:(len(searchQuery)-7)]
	year := searchQuery[(len(searchQuery)-5):(len(searchQuery)-1)]
	url := generateLink(title, useCase)
	body := APIreader(url)

	err := json.Unmarshal(body, &searchResults)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	if len(searchResults.Search) == 0 {
		var searchError SearchError
		fmt.Println("No results found:")		
		json.Unmarshal(body, &searchError)
		fmt.Println(searchError.Error)
	}

	for _, movie := range searchResults.Search {
		if year == movie.Year[0:4] {
			imdbID = movie.IMDbID
			break
		}
	}

	return
}

func getFilmData(imdbID string) (filmData FilmData){// requests the json data by imdbID, returns Film struct
	useCase := "i"
	url := generateLink(imdbID, useCase)
	body := APIreader(url)

	err := json.Unmarshal(body, &filmData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	film = Film{
		ImdbID: 	filmData.ImdbID,
		Title:		filmData.Title,
		Year: 		filmData.Year,
		Director: 	filmData.Director,
		Writer:		filmData.Writer,
		Genre:		filmData.Genre,
		Rating: 	filmData.ImdbRating,
		Plot: 		filmData.Plot,
		Type: 		filmData.Type,
		Runtime:	filmData.Runtime,
		Country:	filmData.Country,
	}

	return
}

func generateLink(query string, useCase string) (url string) {// gets a film title and returns proper omdb url
	switch useCase{
	case "i":		
		url = "http://www.omdbapi.com/?apikey="+apikey+"&i="+query+"&plot=full"
	case "s":
		searchTitle := strings.Replace(strings.ToLower(query), " ", "+", -1)
		url = "http://www.omdbapi.com/?apikey="+apikey+"&s="+searchTitle
	}
	fmt.Println(url)
	return
}

func directoryReader() (filmTitles []string){
	files, _ := ioutil.ReadDir(filmdir)
	for _, file := range files {
		filmTitles = append(filmTitles, file.Name())
	}

	return
}

func APIreader(url string) []uint8{//returns the body from url http request for further parsing
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("error: ", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: received status code", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	return body
}
