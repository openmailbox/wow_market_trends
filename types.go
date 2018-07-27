package main

type AuctionDumpResponse struct {
	Files []AuctionDumpFile `json:"files"`
}

type AuctionDumpFile struct {
	Url          string `json:"url"`
	LastModified int    `json:"lastmodified"`
}

type AuctionRealm struct {
	Name string `json:name`
	Slug string `json:slug`
}

type AuctionFetchResponse struct {
	Realms   []AuctionRealm `json:realms`
	Auctions []Auction      `json:auctions`
}

type Auction struct {
	Auc        int    `json:auc`
	Item       int    `json:item`
	Owner      string `json:owner`
	OwnerRealm string `json:ownerRealm`
	Bid        int    `json:bid`
	Buyout     int    `json:buyout`
	Quantity   int    `json:quantity`
	TimeLeft   string `json:timeLeft`
	Rand       int    `json:rand`
	Seed       int    `json:seed`
	Context    int    `json:context`
}
