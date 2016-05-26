package main

import (
	"encoding/json"
	"flag"
	"math/rand"
	"net/http"
	"time"

	"github.com/fblanco/talks/rtb/bid"
)

var port = flag.String("port", "9090", "http server port")

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	http.HandleFunc("/bid", bidder)
	http.ListenAndServe(":"+*port, nil)
}

func bidder(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	json.NewEncoder(w).Encode(bid.Bid{BidderName: "bidder-" + *port, CPM: rand.Float64() * 10, ElapsedTime: time.Since(t) / time.Millisecond})

}
