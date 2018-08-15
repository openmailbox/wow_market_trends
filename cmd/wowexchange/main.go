package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/openmailbox/wow_market_trends/internal"
)

const dbConnect = "dbname=wow host=/run/postgresql"

var apiKey string

func main() {
	db, err := sql.Open("postgres", dbConnect)
	defer db.Close()
	internal.CheckError(err)

	if len(os.Args) > 1 && strings.Compare(os.Args[1], "serve") == 0 {
		internal.StartServer(db)
	} else {
		apiKey = os.Getenv("BLIZZARD_API_KEY")

		if len(apiKey) == 0 {
			log.Fatal("API key not found. Exiting.")
			os.Exit(1)
		}

		internal.RefreshAuctions(db, apiKey)
		internal.UpdatePeriods(db)
		internal.UpdateItems(db, apiKey)

		log.Println("Done.")
	}
}
