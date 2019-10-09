package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/openmailbox/wow_market_trends/internal"
)

const dbConnect = "dbname=wow host=/run/postgresql"

var clientId string
var clientSecret string

func main() {
	db, err := sql.Open("postgres", dbConnect)
	defer db.Close()
	internal.CheckError(err)

	if len(os.Args) > 1 && strings.Compare(os.Args[1], "serve") == 0 {
		internal.StartServer(db)
	} else {
		clientId = os.Getenv("BLIZZARD_CLIENT_ID")
		clientSecret = os.Getenv("BLIZZARD_CLIENT_SECRET")

		if len(clientId) == 0 || len(clientSecret) == 0 {
			log.Fatal("API client ID and secret are both required. Exiting.")
			os.Exit(1)
		}

		token := internal.FetchToken(clientId, clientSecret)

		internal.RefreshAuctions(db, token)
		internal.UpdatePeriods(db)
		internal.UpdateItems(db, token)

		log.Println("Done.")
	}
}
