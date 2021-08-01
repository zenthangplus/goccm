package goccm

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	cm := New(3)
	for i := 0; i < 10; i++ {
		cm.Wait()
		go func(i int) {
			fmt.Printf("Job %d is running\n", i)
			time.Sleep(200 * time.Millisecond)
			cm.Done()
		}(i)
	}
	cm.WaitAllDone()
}

func TestConcurrency(t *testing.T) {
	var testMaxRunningJobs int
	var counter int32
	incrementTasks := 1000
	maxRunningJobs := 3
	cm := New(maxRunningJobs)
	for i := 0; i < incrementTasks; i++ {
		cm.Wait()
		go func() {
			if cm.RunningCount() > testMaxRunningJobs {
				testMaxRunningJobs = cm.RunningCount()
			}
			atomic.AddInt32(&counter, 1)
			time.Sleep(5 * time.Millisecond)
			cm.Done()
		}()
	}
	cm.WaitAllDone()
	if testMaxRunningJobs > maxRunningJobs {
		t.Errorf("The number of concurrency jobs has exceeded %d. Real result %d.\n", maxRunningJobs, testMaxRunningJobs)
	} else {
		fmt.Printf("max number of goroutines spawned: %d, expected: %d\n", testMaxRunningJobs, maxRunningJobs)
	}
	if counter != int32(incrementTasks) {
		t.Errorf("counter value expected %d. Real result %d.", incrementTasks, counter)
	} else {
		fmt.Printf("tasks executed: %d, expected: %d\n", counter, incrementTasks)
	}
}

// todo: add benchmark tests