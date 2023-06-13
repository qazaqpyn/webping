package workerpool

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Job struct {
	URL string
}

type Result struct {
	URL          string
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

type Pool struct {
	worker      *worker
	workerCount int
	jobs        chan Job
	results     chan Result
	wg          *sync.WaitGroup
	stopped     bool
}

func (r Result) String() string {
	return fmt.Sprintf("URL: %s, Status Code: %d, Response Time: %s, Error: %v", r.URL, r.StatusCode, r.ResponseTime, r.Error)
}

func NewPool(workerCount int, timeout time.Duration, results chan Result) *Pool {
	return &Pool{
		worker:      newWorker(timeout),
		workerCount: workerCount,
		jobs:        make(chan Job),
		results:     results,
		wg:          &sync.WaitGroup{},
	}
}

func (p *Pool) Init() {
	for i := 0; i < p.workerCount; i++ {
		go p.initWorker(i)
	}
}

func (p *Pool) Push(j Job) {
	// if the process stopped we don't want to start new jobs => we just skip and don't enter to next checking phase
	if p.stopped {
		return
	}

	// push created job into the jobs channel
	p.jobs <- j
	p.wg.Add(1)
}

func (p *Pool) Stop() {
	// pool is stopped => no more jobs can be pushed to pool
	p.stopped = true
	// close jobs channel => workers will stop after finishing their jobs
	close(p.jobs)
	// wait for all workers to finish their jobs
	p.wg.Wait()
}

func (p *Pool) initWorker(id int) {
	for job := range p.jobs {
		p.results <- p.worker.process(job)
		p.wg.Done()
	}
	log.Printf("Worker %d finished\n", id)
}
