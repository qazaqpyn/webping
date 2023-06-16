package repository

import (
	"context"

	"github.com/qazaqpyn/webping/domain/audit"
	"github.com/qazaqpyn/webping/domain/results"
	auditMongo "github.com/qazaqpyn/webping/internal/repository/mongo/audit"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type Audit interface {
	Create(ctx context.Context, audit *audit.Audit) error
	Find(ctx context.Context) ([]*audit.MongoAuditGroup, error)
	FindByRequestType(ctx context.Context, requestType int) ([]*audit.Audit, error)
}

type Results interface {
	AddResult(url string, result results.Result)
	GetResult(url string) results.Result
	MaxResponseTime() (url string, result results.Result)
	MinResponseTime() (url string, result results.Result)
}

type Repository struct {
	Audit   Audit
	Results Results
}

func NewRepository(client *mongo.Database) *Repository {
	return &Repository{
		Audit:   auditMongo.NewRepoAudit(client, viper.GetString("mongodb.audit")),
		Results: results.NewResults(),
	}
}
