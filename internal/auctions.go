package internal

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lib/pq"
)

type auctionDumpResponse struct {
	Files []auctionDumpFile `json:"files"`
}

type auctionDumpFile struct {
	URL          string `json:"url"`
	LastModified int    `json:"lastmodified"`
}

type auctionRealm struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type auctionFetchResponse struct {
	Realms   []auctionRealm `json:"realms"`
	Auctions []auction      `json:"auctions"`
}

type auction struct {
	Auc        int    `json:"auc"`
	Item       int    `json:"item"`
	Owner      string `json:"owner"`
	OwnerRealm string `json:"ownerRealm"`
	Bid        int    `json:"bid"`
	Buyout     int    `json:"buyout"`
	Quantity   int    `json:"quantity"`
	TimeLeft   string `json:"timeLeft"`
	Rand       int    `json:"rand"`
	Seed       int    `json:"seed"`
	Context    int    `json:"context"`
}

const realmName = "archimonde"

// RefreshAuctions pulls the latest auctions snapshot from dev.battle.net and stores it in PG
func RefreshAuctions(db *sql.DB, apiKey string) {
	var fileID int

	deleteExisting(db)

	latest := loadLatest(db)
	files := fetchDumps(apiKey)

	for _, file := range files {
		if file.LastModified > latest {
			auctions := fetchAuctions(file)

			createAuctions(db, auctions)

			err := db.QueryRow(`INSERT INTO auction_files (url, last_modified) VALUES ($1, $2) RETURNING id`,
				file.URL, file.LastModified).Scan(&fileID)
			CheckError(err)
		} else {
			log.Printf("Skipping. File too old: %v\n", file.URL)
		}
	}

}

func createAuctions(db *sql.DB, auctions []auction) {
	txn, err := db.Begin()
	CheckError(err)

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

	CheckError(err)

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

		CheckError(err)
	}

	_, err = stmt.Exec()
	CheckError(err)

	err = stmt.Close()
	CheckError(err)

	err = txn.Commit()
	CheckError(err)
}

func deleteExisting(db *sql.DB) {
	log.Println("Deleting previous auction data.")
	_, err := db.Query(`DELETE FROM auctions`)
	CheckError(err)
}

func fetchDumps(apiKey string) []auctionDumpFile {
	log.Println("Fetching auction dump files...")

	resp, err := http.Get("https://us.api.battle.net/wow/auction/data/" + realmName + "?locale=en_US&apikey=" + apiKey)
	CheckError(err)

	var data auctionDumpResponse
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&data)
	CheckError(err)

	log.Printf("Retrieved %v files.\n", len(data.Files))

	return data.Files
}

func fetchAuctions(file auctionDumpFile) []auction {
	log.Printf("Fetching from %v\n", file.URL)

	resp, err := http.Get(file.URL)
	if err != nil {
		log.Printf(err.Error())
		log.Println("\nSkipping file.")
		return make([]auction, 0)
	}

	var data auctionFetchResponse
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&data)
	if err != nil {
		log.Printf(err.Error())
		log.Println("\nSkipping file.")
		return make([]auction, 0)
	}

	return data.Auctions
}

func loadAuctions(db *sql.DB) []auction {
	var auctions []auction
	var itemID, bid, quantity int
	var timeLeft string

	rows, err := db.Query(`SELECT item_id, bid, quantity, time_left FROM auctions`)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&itemID, &bid, &quantity, &timeLeft)
		CheckError(err)
		auctions = append(auctions, auction{Item: itemID, Bid: bid, Quantity: quantity, TimeLeft: timeLeft})
	}

	return auctions
}

func loadLatest(db *sql.DB) int {
	var latest int

	rows, err := db.Query(`SELECT last_modified FROM auction_files ORDER BY last_modified DESC LIMIT 1`)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&latest)
		CheckError(err)
	}

	err = rows.Err()
	CheckError(err)

	log.Printf("Latest file is %v\n", latest)

	return latest
}
