package scheduler

import (
	"time"
)

//DurationWatcher is the interface
type DurationWatcher interface {
	Duration() time.Duration
	Watch() chan time.Duration
}

//Schedule will schedule the execution of the function f, exery Duration(), it will automatically change Duration()
//if a change is made to the desired Duration time
func Schedule(f func(), dn DurationWatcher) {
	c := dn.Watch()
	go func() {
		ticker := time.NewTicker(dn.Duration())
		for {
			select {
			case durationChange := <-c:
				ticker.Stop()
				ticker = time.NewTicker(durationChange)
			case <-ticker.C:
				f()
			}
		}
	}()
}
