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
	Create(ctx context.Context, requestType int, url string, ResponseTime time.Duration) error
	GetAll(ctx context.Context) ([]*audit.MongoAuditGroup, error)
	GetByRequestType(ctx context.Context, requestType int) ([]*audit.MongoAuditResp, error)
}

type Results interface {
	AddResult(url string, result results.Result) error
	GetResult(url string) results.Result
	MaxResponseTime() (string, results.Result)
	MinResponseTime() (string, results.Result)
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
