package main

import (
	"database/sql"
	"log"
	"os"
)

const DbConnnect = "dbname=wow host=/run/postgresql"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func loadItemIds(db *sql.DB) []int {
	var count int

	err := db.QueryRow(`SELECT count(DISTINCT item_id) FROM auctions`).Scan(&count)
	checkError(err)

	ids := make([]int, count)
	i := 0

	rows, err := db.Query(`SELECT DISTINCT item_id FROM auctions`)
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&ids[i])
	}

	err = rows.Err()
	checkError(err)

	return ids
}

func main() {
	api_key := os.Getenv("BLIZZARD_API_KEY")

	if len(api_key) == 0 {
		log.Fatal("API key not found. Exiting.")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", DbConnnect)
	defer db.Close()
	checkError(err)

	refreshAuctions(db, api_key)

	ids := loadItemIds(db)
	log.Printf("Found %v unique items.", len(ids))

	log.Println("Done.")
}
