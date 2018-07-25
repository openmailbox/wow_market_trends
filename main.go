package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const DbConnnect = "user=brandon dbname=wow"

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

func fetchDumps() []AuctionDumpFile {

	resp, err := http.Get("https://us.api.battle.net/wow/auction/data/archimonde?locale=en_US&apikey=" + api_key)
	if err != nil {
		log.Fatalf("Error during fetch: %v", err.Error())
		os.Exit(1)
	}

	var data AuctionDumpResponse
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&data)

	if err != nil {
		log.Fatalf("Error while decoding JSON: %v", err.Error())
		os.Exit(1)
	}

	log.Printf("Retrieved %v files.\n", len(data.Files))

	return data.Files
}

func initDatabase() {
	db, err := sql.Open("postgres", DbConnnect)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
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

	for _, auction := range data.Auctions {
		log.Printf("Inserting %v\n", auction)
	}
}

func main() {
	if len(api_key) == 0 {
		log.Fatal("API key not found. Exiting.")
		os.Exit(1)
	}

	initDatabase()

	files := fetchDumps()

	for _, file := range files {
		fetchAuctions(file)
	}
}
