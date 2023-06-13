package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qazaqpyn/webping/pkg/csvreader"
	"github.com/qazaqpyn/webping/pkg/workerpool"
)

const (
	INTERVAL        = time.Second * 1
	REQUEST_TIMEOUT = time.Second * 2
	WORKER_COUNT    = 10
)

func main() {
	// read urls from csv file
	urls, err := csvreader.ReadCsvFile("./assets/websites.csv")
	if err != nil {
		log.Fatal(err)
		return
	}

	results := make(chan workerpool.Result)
	workerPool := workerpool.NewPool(WORKER_COUNT, REQUEST_TIMEOUT, results)

	// run workers that range over jobs channel
	workerPool.Init()

	// generate jobs
	go generateJobs(workerPool, urls)

	//processResults runnnig in background
	go processResults(results)

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	workerPool.Stop()
}

func generateJobs(workerPool *workerpool.Pool, urls []string) {
	for {
		for _, url := range urls {
			workerPool.Push(workerpool.Job{URL: url})
		}
		time.Sleep(INTERVAL)
	}
}

func processResults(results chan workerpool.Result) {
	go func() {
		for results := range results {
			log.Println(results)
		}
	}()
}
