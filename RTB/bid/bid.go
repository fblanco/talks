package bid

import "time"

//Bid contains the bidder response
type Bid struct {
	CPM         float64
	BidderName  string
	ElapsedTime time.Duration
}

//ProcessedBid contains the bidder response + total round trip time and original url
type ProcessedBid struct {
	OK                   bool
	RoundTripElapsedTime time.Duration
	URL                  string

	Bid
}
