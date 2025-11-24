package scheduler

import "time"

func Scheduler() {
	for {
		time.Sleep(1 * time.Minute)
	}
}
