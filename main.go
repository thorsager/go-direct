package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	version   string
	redirects map[string]string
)

func main() {
	bindAddr := os.Getenv("SERVER_IP")
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	err := json.Unmarshal([]byte(os.Getenv("REDIRECTS")), &redirects)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	for s, t := range redirects {
		log.Printf(" * %s -> %s", s, t)
	}

	http.HandleFunc("/", redirect)

	log.Printf("Starting go-direct(build:%s) on port %s", version, port)
	err = http.ListenAndServe(bindAddr+":"+port, nil)
	if err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	if target, ok := redirects[r.URL.Path]; ok {
		defer log.Printf("Redirected %s -> %s in %v", r.URL.RequestURI(), target, time.Now().Sub(startTime))
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	} else {
		defer log.Printf("No mapping for %s (%v)", r.URL.RequestURI(), time.Now().Sub(startTime))
		http.NotFound(w, r)
	}
}
