package main

import (
	"database/sql"
	"log"
	"os"
	"strings"
)

const dbConnect = "dbname=wow host=/run/postgresql"

var apiKey string

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("postgres", dbConnect)
	defer db.Close()
	checkError(err)

	if len(os.Args) > 1 && strings.Compare(os.Args[1], "serve") == 0 {
		startServer(db)
	} else {
		apiKey = os.Getenv("BLIZZARD_API_KEY")

		if len(apiKey) == 0 {
			log.Fatal("API key not found. Exiting.")
			os.Exit(1)
		}

		refreshAuctions(db, apiKey)
		updatePeriods(db)
		updateItems(db)

		log.Println("Done.")
	}
}
