package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const localAddress = ":8081"

var db *sql.DB

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v / from  %v\n", r.Method, r.RemoteAddr)

	lookupId := 106540 // TODO: change this to a param
	var periods []period

	rows, err := db.Query(`SELECT item_id, high, low, volume, open, close, created_at FROM periods WHERE item_id = $1`, lookupId)
	checkError(err)

	var itemID, high, low, volume, open, close int
	var createdAt time.Time
	var nextPeriod period

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&itemID, &high, &low, &volume, &open, &close, &createdAt)
		checkError(err)

		nextPeriod.ItemID = itemID
		nextPeriod.High = high
		nextPeriod.Low = low
		nextPeriod.Volume = volume
		nextPeriod.Open = open
		nextPeriod.Close = close
		nextPeriod.CreatedAt = createdAt

		periods = append(periods, nextPeriod)
	}

	json.NewEncoder(w).Encode(periods)
	log.Printf("Completed %v %v\n", http.StatusOK, http.StatusText(http.StatusOK))
}

func startServer(database *sql.DB) {
	db = database
	http.HandleFunc("/", handleRequest)
	log.Printf("Listening on %v\n", localAddress)
	log.Fatal(http.ListenAndServe(localAddress, nil))
}
