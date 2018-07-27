package main

import (
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"log"
	"net/http"
)

type AuctionDumpResponse struct {
	Files []AuctionDumpFile `json:"files"`
}

type AuctionDumpFile struct {
	Url          string `json:"url"`
	LastModified int    `json:"lastmodified"`
}

type AuctionRealm struct {
	Name string `json:name`
	Slug string `json:slug`
}

type AuctionFetchResponse struct {
	Realms   []AuctionRealm `json:realms`
	Auctions []Auction      `json:auctions`
}

type Auction struct {
	Auc        int    `json:auc`
	Item       int    `json:item`
	Owner      string `json:owner`
	OwnerRealm string `json:ownerRealm`
	Bid        int    `json:bid`
	Buyout     int    `json:buyout`
	Quantity   int    `json:quantity`
	TimeLeft   string `json:timeLeft`
	Rand       int    `json:rand`
	Seed       int    `json:seed`
	Context    int    `json:context`
}

func refreshAuctions(db *sql.DB, api_key string) {
	var fileId int

	deleteExisting(db)

	latest := loadLatest(db)
	files := fetchDumps(api_key)

	for _, file := range files {
		if file.LastModified > latest {
			auctions := fetchAuctions(file)

			createAuctions(db, auctions)

			err := db.QueryRow(`INSERT INTO auction_files (url, last_modified) VALUES ($1, $2) RETURNING id`,
				file.Url, file.LastModified).Scan(&fileId)
			checkError(err)
		} else {
			log.Printf("Skipping. File too old: %v\n", file.Url)
		}
	}

}

func createAuctions(db *sql.DB, auctions []Auction) {
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

	for _, auction := range auctions {
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

func deleteExisting(db *sql.DB) {
	log.Println("Deleting previous auction data.")
	_, err := db.Query(`DELETE FROM auctions`)
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

func fetchAuctions(file AuctionDumpFile) []Auction {
	log.Printf("Fetching from %v\n", file.Url)

	resp, err := http.Get(file.Url)
	if err != nil {
		log.Printf(err.Error())
		log.Println("\nSkipping file.")
		return make([]Auction, 0)
	}

	var data AuctionFetchResponse
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&data)
	if err != nil {
		log.Printf(err.Error())
		log.Println("\nSkipping file.")
		return make([]Auction, 0)
	}

	return data.Auctions
}

func loadAuctions(db *sql.DB, itemId int) []Auction {
	var auctions []Auction
	var auction_id, bid, buyout, quantity int
	var time_left string

	rows, err := db.Query(`SELECT auction_id, bid, buyout, quantity, time_left FROM auctions WHERE item_id = $1`, itemId)
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&auction_id, &bid, &buyout, &quantity, &time_left)
		checkError(err)
		auctions = append(auctions, Auction{Auc: auction_id, Bid: bid, Buyout: buyout, Quantity: quantity, TimeLeft: time_left})
	}

	return auctions
}

func loadLatest(db *sql.DB) int {
	var latest int

	rows, err := db.Query(`SELECT last_modified FROM auction_files ORDER BY last_modified DESC LIMIT 1`)
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&latest)
		checkError(err)
	}

	err = rows.Err()
	checkError(err)

	log.Printf("Latest file is %v\n", latest)

	return latest
}
