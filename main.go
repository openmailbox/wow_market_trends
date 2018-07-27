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

func main() {
	api_key := os.Getenv("BLIZZARD_API_KEY")

	if len(api_key) == 0 {
		log.Fatal("API key not found. Exiting.")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", DbConnnect)
	defer db.Close()
	checkError(err)

	//refreshAuctions(db, api_key)
	updatePeriods(db)

	log.Println("Done.")
}
