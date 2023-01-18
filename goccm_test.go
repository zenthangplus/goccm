package goccm

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	c := New(3)
	for i := 1; i <= 10; i++ {
		c.Wait()
		go func(i int) {
			fmt.Printf("Job %d is running\n", i)
			time.Sleep(20 * time.Millisecond)
			c.Done()
		}(i)
	}
	c.WaitAllDone()
}

// TestManuallyClose will close after 5 jobs, others should not run
func TestManuallyClose(t *testing.T) {
	executedJobs := make(chan int, 1000)

	c := New(3)
	for i := 1; i <= 1000; i++ {
		jobId := i

		c.Wait()
		go func() {
			executedJobs <- jobId
			fmt.Printf("Executed job id %d\n", jobId)
			time.Sleep(20 * time.Millisecond)
			c.Done()
		}()

		if i == 5 {
			log.Printf("Closing manager")
			c.Close()
			break
		}
	}
	c.WaitAllDone()
}

func TestConcurrency(t *testing.T) {
	var maxRunningJobs = 3
	testMaxRunningJobs := make(chan int32, 100)
	c := New(maxRunningJobs)

	for i := 1; i <= 10; i++ {
		c.Wait()
		go func(i int) {
			fmt.Printf("Current running jobs %d\n", c.RunningCount())
			testMaxRunningJobs <- c.RunningCount()
			time.Sleep(20 * time.Millisecond)
			c.Done()
		}(i)
	}

	c.WaitAllDone()

	for i := 1; i <= 10; i++ {
		observed := <-testMaxRunningJobs

		if observed > int32(maxRunningJobs) {
			t.Errorf("The number of concurrency jobs has exceeded %d. Real result %d.", maxRunningJobs, int(observed))
		}
	}
}
