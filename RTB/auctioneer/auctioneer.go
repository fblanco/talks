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
	port       = flag.String("port", "8000", "http server port")
	bidderList = []string{"http://localhost:9090/bid",
		"http://localhost:9091/bid", "http://localhost:9092/bid",
		"http://localhost:9093/bid", "http://localhost:9094/bid"}
	timeout = 50 * time.Millisecond
)

func main() {
	flag.Parse()
	http.HandleFunc("/ad", auction)
	http.ListenAndServe(":"+*port, nil)
}

func auction(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	ch := make(chan bid.ProcessedBid, len(bidderList))
	callBidders(ch, bidderList)
	bids := collectBids(ch, len(bidderList))
	win := pickWinner(bids)
	if win == nil {
		io.WriteString(w, "nothing returned")
	} else {
		json.NewEncoder(w).Encode(win)
	}
	log.Printf("total time spent:%v", time.Since(t))
}

func pickWinner(bids []bid.ProcessedBid) *bid.ProcessedBid {
	mi, mcpm := -1, -100.0
	for i, b := range bids {
		if !b.OK {
			continue
		}
		if b.CPM > mcpm {
			mi, mcpm = i, b.CPM
		}
	}
	if mi < 0 {
		return nil
	}
	return &bids[mi]
}
func callBidders(ch chan bid.ProcessedBid, bidderList []string) {
	for _, b := range bidderList {
		url := b
		go func() {
			t := time.Now()
			v := getBid(url)
			ok := true
			if v == nil {
				ok = false
				b := bid.Bid{BidderName: "n/a"}
				v = &b
				log.Printf("nothing returned")
			}
			pb := bid.ProcessedBid{OK: ok, RTElapsedTime: time.Since(t), Bid: *v}
			ch <- pb
			log.Printf("%#v", pb)
		}()
	}
}

func collectBids(ch chan bid.ProcessedBid, size int) []bid.ProcessedBid {
	to := time.After(timeout)
	bids := make([]bid.ProcessedBid, 0, size)
	c := 0
Loop:
	for {
		select {
		case pb := <-ch:
			bids = append(bids, pb)
			c++
			if c == size {
				break Loop
			}
		case <-to:
			break Loop
		}
	}

	return bids
}
func getBid(b string) *bid.Bid {
	resp, err := http.Get(b)
	if err != nil {
		log.Printf("error:%s", err.Error())
		return nil
	}
	defer resp.Body.Close()
	var bd bid.Bid
	json.NewDecoder(resp.Body).Decode(&bd)

	return &bd
}
