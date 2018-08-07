package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type item struct {
	ItemID int    `json:"id"`
	Name   string `json:"name"`
}

func updateItems(db *sql.DB) {
	result, err := db.Exec(`INSERT INTO items(item_id) SELECT item_id FROM auctions GROUP BY item_id ON CONFLICT DO NOTHING`)
	checkError(err)

	count, err := result.RowsAffected()
	checkError(err)

	log.Printf("Found %v new item IDs for import.\n", count)

	if count == 0 {
		return
	}

	var nextID int

	rows, err := db.Query(`SELECT item_id FROM items WHERE name IS NULL`)
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&nextID)

		url := fmt.Sprintf("https://us.api.battle.net/wow/item/%v?locale=en_US&apikey=%v", nextID, apiKey)
		resp, err := http.Get(url)
		checkError(err)

		log.Printf("%v", resp)

		var nextItem item

		json.NewDecoder(resp.Body).Decode(&nextItem)

		//db.Exec(`UPDATE items SET name = $1 WHERE item_id = $2`, nextItem.Name, nextItem.ItemID)
		//checkError(err)
	}
}
