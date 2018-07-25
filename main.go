package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/lib/pq"
)

const DbConnnect = "dbname=wow password=brandon"

var db *sql.DB
var latest int
var api_key = os.Getenv("BLIZZARD_API_KEY")

type AuctionDumpFile struct {
	Url          string `json:"url"`
	LastModified int    `json:"lastmodified"`
}

type AuctionDumpResponse struct {
	Files []AuctionDumpFile `json:"files"`
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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func fetchDumps() []AuctionDumpFile {
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

	stmt, err := txn.Prepare(pq.CopyIn("auctions", "auction_id", "item_id", "owner", "owner_realm", "bid", "buyout", "quantity", "time_left", "rand", "seed", "context"))
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

func main() {
	var err error

	if len(api_key) == 0 {
		log.Fatal("API key not found. Exiting.")
		os.Exit(1)
	}

	db, err = sql.Open("postgres", DbConnnect)
	defer db.Close()
	checkError(err)

	rows, err := db.Query(`SELECT last_modified FROM auction_files ORDER BY last_modified DESC LIMIT 1`)
	checkError(err)

	for rows.Next() {
		err = rows.Scan(&latest)
		checkError(err)
	}

	log.Printf("Latest file is %v\n", latest)

	files := fetchDumps()

	for _, file := range files {
		var fileId int

		if file.LastModified > latest {
			fetchAuctions(file)
			err = db.QueryRow(`INSERT INTO auction_files (url, last_modified) VALUES ($1, $2) RETURNING id`, file.Url, file.LastModified).Scan(&fileId)
			checkError(err)
		} else {
			log.Printf("Skipping. File too old: %v\n", file.Url)
		}
	}

	log.Println("Done.")
}
