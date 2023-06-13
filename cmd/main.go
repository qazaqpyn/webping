package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qazaqpyn/webping/pkg/workerpool"
)

const (
	INTERVAL        = time.Second * 1
	REQUEST_TIMEOUT = time.Second * 2
	WORKER_COUNT    = 10
)

var ulrs = []string{
	"https://www.google.com",
	"https://www.facebook.com",
	"https://www.youtube.com",
	"https://www.yahoo.com",
	"https://www.wikipedia.org",
	"https://www.baidu.com",
	"https://www.amazon.com",
	"https://www.qq.com",
	"https://www.google.co.in",
	"https://www.twitter.com",
	"https://www.live.com",
	"https://www.taobao.com",
}

func main() {
	results := make(chan workerpool.Result)
	workerPool := workerpool.NewPool(WORKER_COUNT, REQUEST_TIMEOUT, results)

	// run workers that range over jobs channel
	workerPool.Init()

	// generate jobs
	go generateJobs(workerPool)

	//processResults runnnig in background
	go processResults(results)

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	workerPool.Stop()
}

func generateJobs(workerPool *workerpool.Pool) {
	for {
		for _, url := range ulrs {
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
