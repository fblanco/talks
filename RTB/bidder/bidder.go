package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/fblanco/talks/RTB/bid"
)

var (
	myName *string
	port   *string
)

func init() {
	//seed random number generator
	rand.Seed(time.Now().UnixNano())
	myName = flag.String("name", "bidder-x", "bidder name")
	port = flag.String("port", "8090", "http port number")
	flag.Parse()
}
func main() {
	http.HandleFunc("/bid", biddr)
	fmt.Printf(http.ListenAndServe(":"+*port, nil).Error())
}

func biddr(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	json.NewEncoder(w).Encode(bid.Bid{CPM: rand.Float64() * 10, BidderName: *myName, ElapsedTime: time.Since(t) / time.Millisecond})
}
