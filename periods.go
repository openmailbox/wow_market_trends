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
	log.Println("Updating periods.")

	var itemId, high, low, volume int

	timestamp := time.Now()

	rows, err := db.Query(`SELECT item_id, MAX(bid), MIN(bid), SUM(quantity) FROM auctions GROUP BY item_id`)
	checkError(err)
	defer rows.Close()

	txn, err := db.Begin()
	checkError(err)

	stmt, err := txn.Prepare(pq.CopyIn("periods", "item_id", "high", "low", "volume", "open", "close", "created_at"))
	checkError(err)

	for rows.Next() {
		err = rows.Scan(&itemId, &high, &low, &volume)
		checkError(err)

		period := Period{ItemId: itemId, High: high, Low: low, Volume: volume, CreatedAt: timestamp}

		auctions := loadAuctions(db, period.ItemId)
		period.Close = calculateClose(auctions, period.ItemId)
		period.Open = calculateOpen(db, period.ItemId)

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

func calculateClose(auctions []Auction, itemId int) int {
	var short, medium, long, veryLong []Auction

	for _, auction := range auctions {
		switch auction.TimeLeft {
		case "SHORT":
			short = append(short, auction)
		case "MEDIUM":
			medium = append(medium, auction)
		case "LONG":
			long = append(long, auction)
		case "VERY_LONG":
			veryLong = append(veryLong, auction)
		}
	}

	if len(short) > 0 {
		return bidOverVolume(short)
	} else if len(medium) > 0 {
		return bidOverVolume(medium)
	} else if len(long) > 0 {
		return bidOverVolume(long)
	} else {
		return bidOverVolume(veryLong)
	}
}

func calculateOpen(db *sql.DB, itemId int) int {
	var count, last int

	err := db.QueryRow(`SELECT COUNT(1) FROM periods WHERE item_id = $1`, itemId).Scan(&count)
	checkError(err)

	if count == 0 {
		return 0
	}

	err = db.QueryRow(`SELECT COALESCE(close, 0) FROM periods WHERE item_id = $1 
		ORDER BY created_at DESC LIMIT 1`, itemId).Scan(&last)
	checkError(err)

	return last
}

func bidOverVolume(auctions []Auction) int {
	totalBid := 0
	totalQuantity := 0

	for _, auction := range auctions {
		totalBid += auction.Bid
		totalQuantity += auction.Quantity
	}

	return totalBid / totalQuantity
}
