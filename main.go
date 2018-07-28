package main

import (
	"database/sql"
	"log"
	"os"
)

const dbConnect = "dbname=wow host=/run/postgresql"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	apiKey := os.Getenv("BLIZZARD_API_KEY")

	if len(apiKey) == 0 {
		log.Fatal("API key not found. Exiting.")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", dbConnect)
	defer db.Close()
	checkError(err)

	refreshAuctions(db, apiKey)
	updatePeriods(db)

	log.Println("Done.")
}
