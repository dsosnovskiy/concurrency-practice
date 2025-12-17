package worker

import (
	"fmt"
	"sync"
	"time"
)

func Worker(wg *sync.WaitGroup, id int, jobs <-chan int, results chan<- int) {
	defer wg.Done()

	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(1 * time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}
