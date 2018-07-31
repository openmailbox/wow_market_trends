package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const localAddress = ":8081"

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming request from %v", r.RemoteAddr)

	periods := make([]period, 3)
	encoder := json.NewEncoder(w)

	err := encoder.Encode(periods)
	if err != nil {
		panic(err)
	}
}

func startServer() {
	log.Printf("Listening on %v\n", localAddress)
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(localAddress, nil))
}
