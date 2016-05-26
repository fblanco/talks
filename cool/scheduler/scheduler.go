package scheduler

import (
	"time"
)

//DurationChangeNotifier is the interface
type DurationChangeNotifier interface {
	Duration() time.Duration
	ChangeChannel() chan time.Duration
}

//Schedule will schedule the execution of the function f, exery Duration(), it will automatically change Duration()
//if a change is made to the desired Duration time
func Schedule(f func(), dn DurationChangeNotifier) {
	c := dn.ChangeChannel()
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
