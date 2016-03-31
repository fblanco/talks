package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/fblanco/talks/RTB/bid"
)

var (
	port = flag.String("port", "8090", "http port number")
)

func init() {
	//seed random number generator
	rand.Seed(time.Now().UnixNano())
	port = flag.String("port", "8090", "http port number")
	flag.Parse()
}
func main() {
	http.HandleFunc("/bid", biddr)
	log.Printf(http.ListenAndServe(":"+*port, nil).Error())
}

func biddr(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	json.NewEncoder(w).Encode(bid.Bid{CPM: rand.Float64() * 10, BidderName: "bidder-" + *port, ElapsedTime: time.Since(t) / time.Millisecond})
}
