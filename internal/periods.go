package internal

import (
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
)

type period struct {
	ItemID    int       `json:"item_id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	High      int       `json:"high"`
	Low       int       `json:"low"`
	Volume    int       `json:"volume"`
	Open      int       `json:"open"`
	Close     int       `json:"close"`
	CreatedAt time.Time `json:"created_at"`
	Ask       int       `json:"ask"`
	Average   int       `json:"average"` // 7-day simple moving average
}

// UpdatePeriods takes the current data in the auctions table and aggregates it into a time-period rollup per item
func UpdatePeriods(db *sql.DB) {
	log.Println("Updating periods.")

	var itemID, high, low, volume int

	timestamp := time.Now()
	lastPeriod := loadLastPeriod(db)
	auctions := loadAuctions(db)

	rows, err := db.Query(`SELECT item_id, MAX(bid / quantity), MIN(bid / quantity), SUM(quantity) FROM auctions GROUP BY item_id`)
	CheckError(err)
	defer rows.Close()

	txn, err := db.Begin()
	CheckError(err)

	stmt, err := txn.Prepare(pq.CopyIn("periods", "item_id", "high", "low", "volume", "open", "close", "created_at", "ask"))
	CheckError(err)

	for rows.Next() {
		err = rows.Scan(&itemID, &high, &low, &volume)
		CheckError(err)

		p := period{ItemID: itemID, High: high, Low: low, Volume: volume, CreatedAt: timestamp}

		p.Close, p.Ask = calculateCloseAndAsk(auctions, p.ItemID)
		p.Open = calculateOpen(lastPeriod, p.ItemID)

		_, err = stmt.Exec(p.ItemID, p.High, p.Low, p.Volume, p.Open, p.Close, p.CreatedAt, p.Ask)
		CheckError(err)
	}

	_, err = stmt.Exec()
	CheckError(err)

	err = stmt.Close()
	CheckError(err)

	err = txn.Commit()
	CheckError(err)
}

func calculateCloseAndAsk(auctions []auction, itemID int) (int, int) {
	var short, medium, long, veryLong []auction
	var close, ask int

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
		close = bidOverVolume(short)
		ask = findMinBid(short)
	} else if len(medium) > 0 {
		close = bidOverVolume(medium)
		ask = findMinBid(medium)
	} else if len(long) > 0 {
		close = bidOverVolume(long)
		ask = findMinBid(long)
	} else {
		close = bidOverVolume(veryLong)
		ask = findMinBid(veryLong)
	}

	return close, ask
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

func findMinBid(auctions []auction) int {
	// TODO: use sort pkg
	bid := auctions[0].Bid

	for _, auction := range auctions {
		if auction.Bid < bid {
			bid = auction.Bid
		}
	}

	return bid
}

func loadLastPeriod(db *sql.DB) []period {
	var itemID, close int
	var timestamp time.Time
	var periods []period

	db.QueryRow(`SELECT MAX(created_at) FROM periods`).Scan(&timestamp)

	rows, err := db.Query(`SELECT item_id, close FROM periods WHERE created_at = $1`, timestamp)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&itemID, &close)
		p := period{ItemID: itemID, Close: close}
		periods = append(periods, p)
	}

	return periods
}
