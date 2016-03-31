package bid

import "time"

//Bid contains the bidder response
type Bid struct {
	CPM         float64       `json:"cpm"`
	BidderName  string        `json:"name"`
	ElapsedTime time.Duration `json:"et"`
}

//ProcessedBid contains the bidder response + total round trip time and original url
type ProcessedBid struct {
	OK                   bool          `json:"ok"`
	RoundTripElapsedTime time.Duration `json:"rtet"`
	URL                  string        `json:"url"`

	Bid
}
