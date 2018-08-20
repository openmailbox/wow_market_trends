package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const localAddress = ":8081"
const logFile = "../../log/server.log"

type itemDetails struct {
	Name    string   `json:"name"`
	Icon    string   `json:"icon"`
	Periods []period `json:"periods"`
	Current int      `json:"current"`
}

var db *sql.DB

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "https://wow.open-mailbox.com")
}

func handleDetails(w http.ResponseWriter, r *http.Request) {
	itemIDParam := r.URL.Query()["itemId"]

	if len(itemIDParam) == 0 {
		log.Printf("Completed %v %v\n", http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var details itemDetails

	// TODO: Filter out non-integer itemId param
	lookupID, err := strconv.Atoi(itemIDParam[0])
	CheckError(err)

	details.Periods = fetchItemHistory(lookupID)
	details.Current = details.Periods[0].Ask
	details.Name = details.Periods[0].Name
	details.Icon = details.Periods[0].Icon

	json.NewEncoder(w).Encode(details)
	log.Printf("Completed %v %v\n", http.StatusOK, http.StatusText(http.StatusOK))
}

func handleNameSearch(w http.ResponseWriter, r *http.Request) {
	searchTermParam := r.URL.Query()["search"]

	if len(searchTermParam) == 0 {
		log.Printf("Completed %v %v\n", http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	searchTerm := searchTermParam[0]

	rows, err := db.Query(`SELECT item_id, name FROM items WHERE name ~* $1`, searchTerm)
	CheckError(err)

	var nextItem item
	var items []item
	var itemID int
	var name string

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&itemID, &name)
		CheckError(err)

		nextItem.ItemID = itemID
		nextItem.Name = name

		items = append(items, nextItem)
	}

	json.NewEncoder(w).Encode(items)
	log.Printf("Completed %v %v\n", http.StatusOK, http.StatusText(http.StatusOK))
}

func fetchItemHistory(lookupID int) []period {
	/*
		TODO: New query for daily rollup
		with hourly as (
		select created_at::date as day,
				 created_at,
				   item_id,
				   open,
				   first_value(open) over w as first_open,
				   last_value(close) over w as last_close,
				   high,
				   low,
				   volume,
				   ask
			from periods
			where item_id = 106623
			window w as (partition by periods.created_at::date order by periods.created_at)
			order by periods.created_at desc
		  )
		select day, first_open as open, last_close as close, max(high) as high, min(low) as low, sum(volume) as volume, min(ask) as ask from hourly
		group by day, first_open, last_close
		order by day desc;
	*/
	rows, err := db.Query(`SELECT periods.item_id, name, coalesce(icon, 'noicon'), high, low, volume, open, close, created_at, coalesce(ask, 0) FROM periods
		INNER JOIN items on items.item_id = periods.item_id
		WHERE periods.item_id = $1 ORDER BY periods.id DESC`, lookupID)
	CheckError(err)

	var periods []period
	var itemID, high, low, volume, open, close, ask int
	var name, icon string
	var createdAt time.Time
	var nextPeriod period

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&itemID, &name, &icon, &high, &low, &volume, &open, &close, &createdAt, &ask)
		CheckError(err)

		nextPeriod.ItemID = itemID
		nextPeriod.Name = name
		nextPeriod.Icon = icon
		nextPeriod.High = high
		nextPeriod.Low = low
		nextPeriod.Volume = volume
		nextPeriod.Open = open
		nextPeriod.Close = close
		nextPeriod.CreatedAt = createdAt
		nextPeriod.Ask = ask

		periods = append(periods, nextPeriod)
	}

	return periods
}

func fetchCurentPrice(lookupID int) int {
	var bid int

	rows, err := db.Query(`WITH a AS (select *,
                            CASE  
                            WHEN left(time_left,1) = 'S' THEN 1 
                            WHEN left(time_left,1) = 'M' THEN 2
                            WHEN left(time_left,1) = 'L' THEN 3
                            ELSE 4
                            END AS time_left2
                            FROM auctions
                            WHERE item_id = $1
                            LIMIT 50)
                            SELECT bid FROM a 
                            ORDER BY time_left2 LIMIT 1;`, lookupID)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&bid)
		CheckError(err)
	}

	return bid
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
	CheckError(err)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

	if err != nil {
		log.Fatal("Can't open log file.\n", err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	log.SetOutput(mw)
}

// StartServer starts the HTTP server on the local port
func StartServer(database *sql.DB) {
	db = database

	openLogFile(logFile)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.Handle("/", http.FileServer(http.Dir("../../web/static")))
	http.HandleFunc("/names", handleNameSearch)
	http.HandleFunc("/details", handleDetails)

	fmt.Printf("Logging to %v\n", logFile)
	log.Printf("Listening on %v\n", localAddress)

	log.Fatal(http.ListenAndServe(localAddress, logRequest(http.DefaultServeMux)))
}
