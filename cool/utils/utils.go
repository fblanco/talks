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

// AtomicDuration defines an atomic time.Duration value using sync.atomic and using unit64
type AtomicDuration int64

//Get returns time.Duration
func (d *AtomicDuration) Get() time.Duration { return (time.Duration)(atomic.LoadInt64((*int64)(d))) }

// Set sets time.Duration
func (d *AtomicDuration) Set(t time.Duration) { atomic.StoreInt64((*int64)(d), (int64)(t)) }

//AtomicNotifiableDurationChange is a struct wrapping AtomicDuration + notification channel
type AtomicNotifiableDurationChange struct {
	d AtomicDuration
	c chan time.Duration
}

//Set sets the value
func (t *AtomicNotifiableDurationChange) Set(v time.Duration) {
	ov := t.d.Get()
	t.d.Set(v)
	// if value changes send new value thru channel
	if ov != 0 && ov != v {
		t.c <- v
	}
}

//Duration gets the value to implement DurationChangeNotifier interface
func (t *AtomicNotifiableDurationChange) Duration() time.Duration {
	return t.d.Get()
}

//ChangeChannel returns the channel where changes notification will be pushed
func (t *AtomicNotifiableDurationChange) ChangeChannel() chan time.Duration {
	if t.c == nil {
		t.c = make(chan time.Duration, 1)
	}
	return t.c
}
