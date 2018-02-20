package main

import (
	"net/http"
	"log"
	"net"
	"time"

	"github.com/gorilla/mux"
	"github.com/coreos/go-systemd/daemon"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.NotFoundHandler = http.HandlerFunc(handle404)

	// Serve files from "static" and "templates" directories
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.PathPrefix("/static/").Handler(s)
	t := http.StripPrefix("/templates/", http.FileServer(http.Dir("templates")))
	r.PathPrefix("/templates/").Handler(t)

	// Valid Routes
	routes := []string{
			"/",
			"/about",
			"/projects",
			"/blog",
			"/projects/{projectname}",
			"/blog/{postname}",
			"/404",
			}

	for _, value := range routes {
			r.HandleFunc(value, pageHandler)
	}

	// Listen
	l, err := net.Listen("tcp", "127.0.0.1:5000")
	if err != nil {
		log.Println(err)
	}

	// Tell systemd website operational.
	daemon.SdNotify(false, "READY=1")

	// Heartbeat
	go func() {
    interval, err := daemon.SdWatchdogEnabled(false)
    if err != nil || interval == 0 {
        return
    }
		for {
	    _, err := http.Get("http://127.0.0.1:5000")
	    if err == nil {
	        daemon.SdNotify(false, "WATCHDOG=1")
	    }
	    time.Sleep(interval / 3)
		}
	}()

	// Serve
  http.Serve(l, r)
}
