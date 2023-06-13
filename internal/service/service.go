package service

import (
	"context"
	"time"

	"github.com/qazaqpyn/webping/domain/audit"
	"github.com/qazaqpyn/webping/domain/results"
	"github.com/qazaqpyn/webping/internal/repository"
	auditService "github.com/qazaqpyn/webping/internal/service/audit"
	authService "github.com/qazaqpyn/webping/internal/service/auth"
	resultsService "github.com/qazaqpyn/webping/internal/service/results"
)

type Audit interface {
	Create(ctx context.Context, requestType int, url string, ResponseTime time.Duration) (*audit.Audit, error)
	GetAll(ctx context.Context) ([]*audit.Audit, error)
	GetByRequestType(ctx context.Context, requestType int) ([]*audit.Audit, error)
}

type Results interface {
	AddResult(url string, result results.Result) error
	GetResult(url string) results.Result
	MaxResponseTime() (string, results.Result)
	MinResponseTime() (string, results.Result)
	CheckWebExist(url string) bool
}

type Auth interface {
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

type Service struct {
	Audit   Audit
	Results Results
	Auth    Auth
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Audit:   auditService.NewAuditService(repos.Audit),
		Results: resultsService.NewResultService(repos.Results),
		Auth:    authService.NewAuthService(),
	}
}
