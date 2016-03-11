package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/fblanco/talks/RTB/bid"
)

var (
	bidderList []string
	timeout    time.Duration
	port       *string
)

func init() {
	bidderList = []string{"http://localhost:8085/bid", "http://localhost:8086/bid", "http://localhost:8087/bid", "http://localhost:8088/bid"}
	timeout = 50 * time.Millisecond
	port = flag.String("port", "8080", "http port number")
	flag.Parse()
}

func main() {
	http.HandleFunc("/ad", auction)
	fmt.Printf(http.ListenAndServe(":"+*port, nil).Error())
}

func auction(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	ch := make(chan bid.ProcessedBid, len(bidderList)+1)
	callBidders(ch, bidderList)
	bids := collectBids(ch, len(bidderList))
	pickWinner(bids)
	fmt.Printf("total et %v\n", time.Since(t))
}

func callBidders(ch chan bid.ProcessedBid, bidderList []string) {
	for _, b := range bidderList {
		url := b
		go func() {
			t := time.Now()
			bd := getBid(url)
			ok := true
			if bd == nil {
				ok = false
				bd = &bid.Bid{BidderName: "n/a"}
			}
			pbid := bid.ProcessedBid{ok, time.Since(t) / time.Millisecond, url, *bd}
			ch <- pbid
		}()
	}
}

func getBid(url string) *bid.Bid {
	httpClient := http.Client{Timeout: timeout * 10}
	resp, err := httpClient.Get(url)
	if err != nil {
		fmt.Printf("error %s\n", err)
		return nil
	}

	var b bid.Bid
	json.NewDecoder(resp.Body).Decode(&b)
	resp.Body.Close()
	return &b
}

func collectBids(ch chan bid.ProcessedBid, size int) []bid.ProcessedBid {
	var bids []bid.ProcessedBid
	notdone := true
	c := 0
	to := time.After(timeout)
	for notdone {
		select {
		case bid := <-ch:
			bids = append(bids, bid)
			c++
			if c == size {
				notdone = false
			}
		case <-to:
			notdone = false
		}
	}
	return bids
}

func pickWinner(bids []bid.ProcessedBid) {
	max := -10.0
	wi := -1
	for i, b := range bids {
		fmt.Printf("%#v\n", b)
		if !b.OK {
			continue
		}
		if b.CPM > max {
			wi = i
			max = b.CPM
		}
	}
	if wi < 0 {
		fmt.Println("nothing came back")
	} else {
		fmt.Printf("the winner is: %#v\n", bids[wi])
	}
}
