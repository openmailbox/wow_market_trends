package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const localAddress = ":8081"
const logFile = "log/server.log"

var db *sql.DB

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func handleNameSearch(w http.ResponseWriter, r *http.Request) {
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

func handleHistory(w http.ResponseWriter, r *http.Request) {
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

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v from %v\n", r.Method, r.URL, r.RemoteAddr)
		handler.ServeHTTP(w, r)
	})
}

func openLogFile(fileName string) {
	if len(fileName) == 0 {
		return
	}

	path, err := filepath.Abs(fileName)
	checkError(err)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

	if err != nil {
		log.Fatal("Can't open log file.", err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	log.SetOutput(mw)
}

func startServer(database *sql.DB) {
	db = database

	openLogFile(logFile)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/history", handleHistory)
	http.HandleFunc("/names", handleNameSearch)

	fmt.Printf("Logging to %v\n", logFile)
	log.Printf("Listening on %v\n", localAddress)

	log.Fatal(http.ListenAndServe(localAddress, logRequest(http.DefaultServeMux)))
}
