package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	// Serve files from "static" and "templates" directories
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.PathPrefix("/static/").Handler(s)
	t := http.StripPrefix("/templates/", http.FileServer(http.Dir("templates")))
	r.PathPrefix("/templates/").Handler(t)

	// Valid Routes
	routes := []string{
			"/",
			"/projects",
			"/blog",
			"/projects/{projectname}",
			"/blog/{postname}"}

	for _, value := range routes {
			r.HandleFunc(value, pageHandler)
	}

  http.ListenAndServe("localhost:5000", r)
}
