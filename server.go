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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func handleNameSearch(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v / from  %v with params: %v\n", r.Method, r.RemoteAddr, r.URL.Query())

	searchTermParam := r.URL.Query()["search"]

	if len(searchTermParam) == 0 {
		log.Printf("Completed %v %v\n", http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	enableCors(&w)

	searchTerm := searchTermParam[0]

	rows, err := db.Query(`SELECT item_id, name FROM items WHERE name ~* $1`, searchTerm)
	checkError(err)

	var nextItem item
	var items []item
	var itemID int
	var name string

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&itemID, &name)
		checkError(err)

		nextItem.ItemID = itemID
		nextItem.Name = name

		items = append(items, nextItem)
	}

	json.NewEncoder(w).Encode(items)
	log.Printf("Completed %v %v\n", http.StatusOK, http.StatusText(http.StatusOK))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v / from  %v with params: %v\n", r.Method, r.RemoteAddr, r.URL.Query())

	itemIDParam := r.URL.Query()["itemId"]

	if len(itemIDParam) == 0 {
		log.Printf("Completed %v %v\n", http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	enableCors(&w)

	lookupID := itemIDParam[0]

	rows, err := db.Query(`SELECT periods.item_id, name, high, low, volume, open, close, created_at FROM periods 
		INNER JOIN items on items.item_id = periods.item_id
		WHERE periods.item_id = $1 ORDER BY periods.id DESC`, lookupID)
	checkError(err)

	var periods []period
	var itemID, high, low, volume, open, close int
	var name string
	var createdAt time.Time
	var nextPeriod period

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&itemID, &name, &high, &low, &volume, &open, &close, &createdAt)
		checkError(err)

		nextPeriod.ItemID = itemID
		nextPeriod.Name = name
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
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/history", handleRequest)
	http.HandleFunc("/names", handleNameSearch)
	log.Printf("Listening on %v\n", localAddress)
	log.Fatal(http.ListenAndServe(localAddress, nil))
}
