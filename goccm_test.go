package goccm

import (
	"fmt"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	c := New(3)
	for i := 1; i <= 10; i++ {
		c.Wait()
		go func(i int) {
			fmt.Printf("Job %d is running\n", i)
			time.Sleep(2 * time.Second)
			c.Done()
		}(i)
	}
	c.WaitAllDone()
}

func TestManuallyClose(t *testing.T) {
	executedJobs := 0
	c := New(3)
	for i := 1; i <= 1000; i++ {
		c.Wait()
		go func() {
			executedJobs++
			fmt.Printf("Executed jobs %d\n", executedJobs)
			time.Sleep(2 * time.Second)
			c.Done()
		}()
		if i == 5 {
			c.Close()
			break
		}
	}
	c.WaitAllDone()
}

func TestConcurrency(t *testing.T) {
	var maxRunningJobs = 3
	var testMaxRunningJobs int32
	c := New(maxRunningJobs)
	for i := 1; i <= 10; i++ {
		c.Wait()
		go func(i int) {
			fmt.Printf("Current running jobs %d\n", c.RunningCount())
			if c.RunningCount() > testMaxRunningJobs {
				testMaxRunningJobs = c.RunningCount()
			}
			time.Sleep(2 * time.Second)
			c.Done()
		}(i)
	}
	c.WaitAllDone()
	if testMaxRunningJobs > int32(maxRunningJobs) {
		t.Errorf("The number of concurrency jobs has exceeded %d. Real result %d.", maxRunningJobs, testMaxRunningJobs)
	}
}
