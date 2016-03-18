package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"time"

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
	log.Printf(http.ListenAndServe(":"+*port, nil).Error())
}

func auction(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	ch := make(chan bid.ProcessedBid, len(bidderList)+1)
	callBidders(ch, bidderList)
	bids := collectBids(ch, len(bidderList))
	if win := pickWinner(bids); win == nil {
		io.WriteString(w, "nothing")
	} else {
		json.NewEncoder(w).Encode(*win)
	}
	log.Printf("total et %v\n", time.Since(t))
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
		log.Printf("error %s\n", err)
		return nil
	}

	var b bid.Bid
	json.NewDecoder(resp.Body).Decode(&b)
	resp.Body.Close()
	return &b
}

func collectBids(ch chan bid.ProcessedBid, size int) []bid.ProcessedBid {
	bids := make([]bid.ProcessedBid, 0, size)
	to := time.After(timeout)
Loop:
	for {
		select {
		case bid := <-ch:
			bids = append(bids, bid)
			if len(bids) == size {
				break Loop
			}
		case <-to:
			break Loop
		}
	}
	return bids
}

func pickWinner(bids []bid.ProcessedBid) *bid.ProcessedBid {
	max, wi := -10.0, -1
	for i, b := range bids {
		log.Printf("checking bid: %#v\n", b)
		if !b.OK {
			continue
		}
		if b.CPM > max {
			max, wi = b.CPM, i
		}
	}
	if wi < 0 {
		return nil
	}
	return &bids[wi]
}
