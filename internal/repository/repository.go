package repository

import (
	"context"

	"github.com/qazaqpyn/webping/domain/audit"
	"github.com/qazaqpyn/webping/domain/results"
	auditMongo "github.com/qazaqpyn/webping/internal/repository/mongo/audit"
	"go.mongodb.org/mongo-driver/mongo"
)

type Audit interface {
	Create(ctx context.Context, audit *audit.Audit) error
	Find(ctx context.Context) ([]*audit.Audit, error)
	FindByRequestType(ctx context.Context, requestType int) ([]*audit.Audit, error)
}

type Results interface {
	AddResult(url string, result results.Result)
	GetResult(url string) results.Result
	MaxResponseTime() (url string, result results.Result)
	MinResponseTime() (url string, result results.Result)
	CheckWebExist(url string) bool
}

type Repository struct {
	Audit   Audit
	Results Results
}

func NewRepository(client mongo.Client) *Repository {
	return &Repository{
		Audit:   auditMongo.NewRepoAudit(client),
		Results: results.NewResults(),
	}
}
