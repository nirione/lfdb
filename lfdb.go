package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
	
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	dat, _ := os.ReadFile("./test")
	fmt.Fprintf(w, string(dat))

}
