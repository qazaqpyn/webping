package results

import (
	"github.com/qazaqpyn/webping/domain/results"
	"github.com/qazaqpyn/webping/internal/repository"
)

type resultService struct {
	resultRepo repository.Results
}

func NewResultService(repo repository.Results) *resultService {
	return &resultService{resultRepo: repo}
}

func (s *resultService) AddResult(url string, result results.Result) error {
	s.resultRepo.AddResult(url, result)
	return nil
}

func (s *resultService) CheckWebExist(url string) bool {
	return s.resultRepo.CheckWebExist(url)
}

func (s *resultService) GetResult(url string) results.Result {
	return s.resultRepo.GetResult(url)
}

func (s *resultService) MaxResponseTime() (string, results.Result) {
	return s.resultRepo.MaxResponseTime()
}

func (s *resultService) MinResponseTime() (string, results.Result) {
	return s.resultRepo.MinResponseTime()
}
