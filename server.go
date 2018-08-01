package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const localAddress = ":8081"

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v / from  %v\n", r.Method, r.RemoteAddr)

	periods := make([]period, 3)

	json.NewEncoder(w).Encode(periods)
	log.Printf("Completed %v %v\n", http.StatusOK, http.StatusText(http.StatusOK))
}

func startServer() {
	http.HandleFunc("/", handleRequest)
	log.Printf("Listening on %v\n", localAddress)
	log.Fatal(http.ListenAndServe(localAddress, nil))
}
