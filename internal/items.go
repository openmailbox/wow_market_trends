package internal

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
	Icon   string `json:"icon"`
}

// UpdateItems uses the battle.net API to fetch item details for item IDs in the periods table
func UpdateItems(db *sql.DB, token string) {
	result, err := db.Exec(`INSERT INTO items(item_id) SELECT item_id FROM auctions GROUP BY item_id ON CONFLICT DO NOTHING`)
	CheckError(err)

	count, err := result.RowsAffected()
	CheckError(err)

	log.Printf("Found %v new item IDs for import.\n", count)

	var nextID int
	var nextItem item

	db.QueryRow(`SELECT COUNT(1) FROM items WHERE name IS NULL OR icon IS NULL`).Scan(&count)
	log.Printf("Updating item data for %v items.", count)

	rows, err := db.Query(`SELECT item_id FROM items WHERE name IS NULL OR icon IS NULL`)
	CheckError(err)

	stmt, err := db.Prepare(`UPDATE items SET name = $1, icon = $2 WHERE item_id = $3`)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&nextID)
		CheckError(err)

		url := fmt.Sprintf("https://us.api.blizzard.com/wow/item/%v?locale=en_US&access_token=%v", nextID, token)
		resp, err := http.Get(url)
		CheckError(err)

		json.NewDecoder(resp.Body).Decode(&nextItem)

		_, err = stmt.Exec(nextItem.Name, nextItem.Icon, nextItem.ItemID)
		CheckError(err)
	}

	_, err = stmt.Exec()

	err = stmt.Close()
	CheckError(err)
}
