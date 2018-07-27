package main

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"time"
)

type Period struct {
	ItemId    int
	High      int
	Low       int
	Volume    int
	Open      int
	Close     int
	CreatedAt time.Time
}

func updatePeriods(db *sql.DB) {
	timestamp := time.Now()
	ids := loadItemIds(db)
	log.Printf("Found %v unique items.", len(ids))

	txn, err := db.Begin()
	checkError(err)

	stmt, err := txn.Prepare(pq.CopyIn("periods", "item_id", "high", "low", "volume", "open", "close", "created_at"))
	checkError(err)

	for _, id := range ids {
		period := Period{ItemId: id}

		period.High = calculateHigh(db, period.ItemId)
		period.Low = calculateLow(db, period.ItemId)
		period.Volume = calculateVolume(db, period.ItemId)
		period.CreatedAt = timestamp

		_, err = stmt.Exec(period.ItemId, period.High, period.Low, period.Volume, period.Open, period.Close, period.CreatedAt)
		checkError(err)
	}

	_, err = stmt.Exec()
	checkError(err)

	err = stmt.Close()
	checkError(err)

	err = txn.Commit()
	checkError(err)
}

func calculateHigh(db *sql.DB, itemId int) int {
	var max int

	err := db.QueryRow(`SELECT MAX(bid) FROM auctions WHERE item_id = $1`, itemId).Scan(&max)
	checkError(err)

	return max
}

func calculateLow(db *sql.DB, itemId int) int {
	var min int

	err := db.QueryRow(`SELECT MIN(bid) FROM auctions WHERE item_id = $1`, itemId).Scan(&min)
	checkError(err)

	return min
}

func calculateVolume(db *sql.DB, itemId int) int {
	var volume int

	err := db.QueryRow(`SELECT SUM(quantity) FROM auctions WHERE item_id = $1`, itemId).Scan(&volume)
	checkError(err)

	return volume
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
		i += 1
	}

	err = rows.Err()
	checkError(err)

	return ids
}
