package results

import (
	"sync"
	"time"
)

type Result struct {
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

// type results map with url and result
type Results struct {
	Results map[string]Result
	mutex   sync.Mutex
}

func NewResult(StatusCode int, ResponseTime time.Duration, Error error) Result {
	return Result{
		StatusCode:   StatusCode,
		ResponseTime: ResponseTime,
		Error:        Error,
	}
}

func NewResults() *Results {
	return &Results{
		Results: make(map[string]Result),
	}
}

func (r *Results) AddResult(url string, result Result) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.Results[url] = result
}

func (r *Results) GetResult(url string) Result {
	for {
		if _, ok := r.Results[url]; ok {
			return r.Results[url]
		}
	}
}

func (r *Results) CheckWebExist(url string) bool {
	// list
	return true
}

func (r *Results) MaxResponseTime() (url string, result Result) {
	for k, v := range r.Results {
		if result.ResponseTime < v.ResponseTime {
			url = k
			result = v
		}
	}
	return
}

func (r *Results) MinResponseTime() (url string, result Result) {
	for k, v := range r.Results {
		if result.ResponseTime > v.ResponseTime {
			url = k
			result = v
		}
	}
	return
}
