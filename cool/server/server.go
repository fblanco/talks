package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fblanco/talks/cool/utils"
)

func signalsHandler(wg *sync.WaitGroup, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

	// Block until a signal is received.
	s := <-c
	//changing healthCheck status return code, give time to clear active requests before exiting
	stop.Set()

	log.Printf("Signal '%s' received, letting requests finish", s)

	// wait for all pending request to finish, but don't wait more than "timeout"
	t := utils.WaitTimeout(wg, timeout)
	if t {
		log.Println("finished waiting because of time out, proceeding with cleanup")
	} else {
		log.Println("finished waiting all pending requests are done, proceeding with cleanup")
	}
	/* clean up here
	   (close db, files etc...)
	*/

	// signal assigned to restart behavior
	if s == syscall.SIGUSR2 {
		log.Println("restarting itself!")
		restart()
		return
	}

	log.Println("bye!")
	os.Exit(1)
}

func restart() {
	syscall.Exec(os.Args[0], os.Args, os.Environ())
}

var port = flag.String("port", "8001", "http server port")
var wg sync.WaitGroup

var stop utils.AtomicBool

func init() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Println("starting up v1")
	go signalsHandler(&wg, 10*time.Second)
	http.HandleFunc("/do", checkInOut(logIt(doSomething)))
	http.HandleFunc("/health", logIt(healthCheck))

	log.Printf(http.ListenAndServe(":"+*port, nil).Error())
}

func checkInOut(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		defer wg.Done()
		f(w, r)
	}
}

func logIt(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		f(w, r)
		log.Printf("m:%s, r:%s, ua:%s, et:%v", r.Method, r.RequestURI, r.UserAgent(), time.Since(t))
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if stop.Get() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	io.WriteString(w, "ok")
}

func doSomething(w http.ResponseWriter, r *http.Request) {
	st := time.Duration(rand.Intn(20)) * time.Second
	time.Sleep(st)
	fmt.Fprintf(w, "took a nap for %v\n", st)
}
