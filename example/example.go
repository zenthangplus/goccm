package main

import (
	"fmt"
	"goccm"
	"time"
)

func main() {
	fmt.Println("Start")
	c := goccm.New(3)
	numberJobs := 10
	for i := 1; i <= numberJobs; i++ {
		c.Wait()
		go func(i int) {
			fmt.Printf("+++ Job %d is running\n", i)
			time.Sleep(3 * time.Second)
			fmt.Printf("--- Job %d is stopped\n", i)
			c.Done()
		}(i)
	}
	c.Close()
	c.WaitAllDone()
	fmt.Println("End")
}
