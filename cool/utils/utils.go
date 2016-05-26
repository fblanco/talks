package utils

import (
	"sync"
	"sync/atomic"
	"time"
)

// WaitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

// AtomicBool defines an atomic boolean value using sync.atomic and using int32
type AtomicBool int32

//Get returns true or false
func (b *AtomicBool) Get() bool { return atomic.LoadInt32((*int32)(b)) != 0 }

// Set  it true (1, !=0)
func (b *AtomicBool) Set() { atomic.StoreInt32((*int32)(b), 1) }
