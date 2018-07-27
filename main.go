package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/lib/pq"
)

const DbConnnect = "dbname=wow host=/run/postgresql"

var db *sql.DB

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func deleteExisting() {
	log.Println("Deleting previous auction data.")
	_, err := db.Query(`DELETE FROM auctions`)
	checkError(err)
}

func fetchAuctions(file AuctionDumpFile) {
	log.Printf("Fetching from %v\n", file.Url)

	resp, err := http.Get(file.Url)
	if err != nil {
		log.Printf(err.Error())
		log.Println("\nSkipping file.")
		return
	}

	var data AuctionFetchResponse
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&data)
	if err != nil {
		log.Printf(err.Error())
		log.Println("\nSkipping file.")
		return
	}

	txn, err := db.Begin()
	checkError(err)

	stmt, err := txn.Prepare(pq.CopyIn(
		"auctions",
		"auction_id",
		"item_id",
		"owner",
		"owner_realm",
		"bid",
		"buyout",
		"quantity",
		"time_left",
		"rand",
		"seed",
		"context"))

	checkError(err)

	for _, auction := range data.Auctions {
		_, err = stmt.Exec(
			auction.Auc,
			auction.Item,
			auction.Owner,
			auction.OwnerRealm,
			auction.Bid,
			auction.Buyout,
			auction.Quantity,
			auction.TimeLeft,
			auction.Rand,
			auction.Seed,
			auction.Context)

		checkError(err)
	}

	_, err = stmt.Exec()
	checkError(err)

	err = stmt.Close()
	checkError(err)

	err = txn.Commit()
	checkError(err)
}

func fetchDumps(api_key string) []AuctionDumpFile {
	log.Println("Fetching auction dump files...")

	resp, err := http.Get("https://us.api.battle.net/wow/auction/data/archimonde?locale=en_US&apikey=" + api_key)
	checkError(err)

	var data AuctionDumpResponse
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&data)
	checkError(err)

	log.Printf("Retrieved %v files.\n", len(data.Files))

	return data.Files
}

func loadLatest() int {
	var latest int

	rows, err := db.Query(`SELECT last_modified FROM auction_files ORDER BY last_modified DESC LIMIT 1`)
	checkError(err)

	for rows.Next() {
		err = rows.Scan(&latest)
		checkError(err)
	}

	log.Printf("Latest file is %v\n", latest)

	return latest
}

func main() {
	var err error
	var fileId int

	api_key := os.Getenv("BLIZZARD_API_KEY")

	if len(api_key) == 0 {
		log.Fatal("API key not found. Exiting.")
		os.Exit(1)
	}

	db, err = sql.Open("postgres", DbConnnect)
	defer db.Close()
	checkError(err)

	deleteExisting()

	latest := loadLatest()
	files := fetchDumps(api_key)

	for _, file := range files {
		if file.LastModified > latest {
			fetchAuctions(file)

			err = db.QueryRow(`INSERT INTO auction_files (url, last_modified) VALUES ($1, $2) RETURNING id`,
				file.Url, file.LastModified).Scan(&fileId)
			checkError(err)
		} else {
			log.Printf("Skipping. File too old: %v\n", file.Url)
		}
	}

	log.Println("Done.")
}
