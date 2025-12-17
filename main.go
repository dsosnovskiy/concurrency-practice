package main

import (
	"concurrency/miner"
	"concurrency/postman"
	"concurrency/worker"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var coal atomic.Int64

	mu := sync.Mutex{}
	var mails []string

	minerContext, minerCancel := context.WithCancel(context.Background())
	postmanContext, postmanCancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("--->>> Рабочий день шахтеров окончен!")
		minerCancel()
	}()

	go func() {
		time.Sleep(6 * time.Second)
		fmt.Println("--->>> Рабочий день почтальонов окончен!")
		postmanCancel()
	}()

	coalTransferPoint := miner.MinerPool(minerContext, 3)
	mailTransferPoint := postman.PostmanPool(postmanContext, 3)

	initTime := time.Now()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for v := range coalTransferPoint {
			coal.Add(int64(v))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for v := range mailTransferPoint {
			mu.Lock()
			mails = append(mails, v)
			mu.Unlock()
		}
	}()

	wg.Wait()

	fmt.Println("Суммарно добытый уголь:", coal.Load())

	mu.Lock()
	fmt.Println("Суммарное количество полученных писем:", len(mails))
	mu.Unlock()

	fmt.Println("Затраченное время:", time.Since(initTime))

	// --------------- Классический Worker Pool ---------------

	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker.Worker(wg, w, jobs, results)
	}

	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		_ = result
	}
}
