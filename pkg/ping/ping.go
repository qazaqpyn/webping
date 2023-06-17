package ping

import (
	"time"

	"github.com/qazaqpyn/webping/domain/results"
	"github.com/qazaqpyn/webping/internal/service"
	"github.com/qazaqpyn/webping/pkg/workerpool"
)

const (
	INTERVAL        = time.Second * 1
	REQUEST_TIMEOUT = time.Second * 10
	WORKER_COUNT    = 10
)

type Ping struct {
	resultsService service.Results
	results        chan workerpool.Result
	workerPool     *workerpool.Pool
}

func NewPingHandler(resultsService service.Results) *Ping {
	results := make(chan workerpool.Result)
	workerPool := workerpool.NewPool(WORKER_COUNT, REQUEST_TIMEOUT, results)

	return &Ping{
		resultsService: resultsService,
		results:        results,
		workerPool:     workerPool,
	}
}

func (h *Ping) StartPing() {
	h.workerPool.Init()
}

func (h *Ping) GenerateJobs(urls []string) {
	for {
		for _, url := range urls {
			h.workerPool.Push(workerpool.Job{URL: url})
		}
		time.Sleep(INTERVAL)
	}
}

func (h *Ping) ProcessResults() {
	go func() {
		for result := range h.results {
			var response results.Result
			response.StatusCode = result.StatusCode
			response.ResponseTime = result.ResponseTime
			response.Error = result.Error

			h.resultsService.AddResult(result.URL, response)
		}
	}()
}
