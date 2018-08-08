package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
)

type period struct {
	ItemID    int       `json:"item_id"`
	High      int       `json:"high"`
	Low       int       `json:"low"`
	Volume    int       `json:"volume"`
	Open      int       `json:"open"`
	Close     int       `json:"close"`
	CreatedAt time.Time `json:"created_at"`
}

func updatePeriods(db *sql.DB) {
	log.Println("Updating periods.")

	var itemID, high, low, volume int

	timestamp := time.Now()
	lastPeriod := loadLastPeriod(db)
	auctions := loadAuctions(db)

	rows, err := db.Query(`SELECT item_id, MAX(bid), MIN(bid), SUM(quantity) FROM auctions GROUP BY item_id`)
	checkError(err)
	defer rows.Close()

	txn, err := db.Begin()
	checkError(err)

	stmt, err := txn.Prepare(pq.CopyIn("periods", "item_id", "high", "low", "volume", "open", "close", "created_at"))
	checkError(err)

	for rows.Next() {
		err = rows.Scan(&itemID, &high, &low, &volume)
		checkError(err)

		p := period{ItemID: itemID, High: high, Low: low, Volume: volume, CreatedAt: timestamp}

		p.Close = calculateClose(auctions, p.ItemID)
		p.Open = calculateOpen(lastPeriod, p.ItemID)

		_, err = stmt.Exec(p.ItemID, p.High, p.Low, p.Volume, p.Open, p.Close, p.CreatedAt)
		checkError(err)
	}

	_, err = stmt.Exec()
	checkError(err)

	err = stmt.Close()
	checkError(err)

	err = txn.Commit()
	checkError(err)
}

func calculateClose(auctions []auction, itemID int) int {
	var short, medium, long, veryLong []auction

	for _, auction := range auctions {
		if auction.Item != itemID {
			continue
		}

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

func calculateOpen(lastPeriod []period, itemID int) int {
	last := 0

	for _, period := range lastPeriod {
		if period.ItemID == itemID {
			last = period.Close
			break
		}
	}

	// TODO: Either go back further or use the current price
	return last
}

func bidOverVolume(auctions []auction) int {
	totalBid := 0
	totalQuantity := 0

	for _, auction := range auctions {
		totalBid += auction.Bid
		totalQuantity += auction.Quantity
	}

	return totalBid / totalQuantity
}

func loadLastPeriod(db *sql.DB) []period {
	var itemID, close int
	var timestamp time.Time
	var periods []period

	db.QueryRow(`SELECT MAX(created_at) FROM periods`).Scan(&timestamp)

	rows, err := db.Query(`SELECT item_id, close FROM periods WHERE created_at = $1`, timestamp)
	checkError(err)

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&itemID, &close)
		p := period{ItemID: itemID, Close: close}
		periods = append(periods, p)
	}

	return periods
}
