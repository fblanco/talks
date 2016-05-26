package bid

import "time"

// Bid is a struct containing bid info
type Bid struct {
	BidderName  string `json:"name"`
	CPM         float64
	ElapsedTime time.Duration
}

// ProcessedBid is a struct containing bid info + auctioneer info
type ProcessedBid struct {
	OK            bool
	RTElapsedTime time.Duration

	Bid
}
