package goccm

import (
	"strconv"
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
	var testMaxRunningJobs int
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
	if testMaxRunningJobs > maxRunningJobs {
		t.Errorf("The number of concurrency jobs has exceeded %d. Real result %d.", maxRunningJobs, testMaxRunningJobs)
	}
}

func BenchmarkConcurrency(b *testing.B) {
	for i := 3; i <= 30; i++ {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			c := New(i)
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				c.Wait()
				go c.Done()
			}
			c.WaitAllDone()
		})
	}
}
