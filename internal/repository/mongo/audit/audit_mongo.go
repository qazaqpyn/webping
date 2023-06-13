package audit

import (
	"context"
	"time"

	"github.com/qazaqpyn/webping/domain/audit"
	"github.com/qazaqpyn/webping/domain/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepoAudit struct {
	client mongo.Client
	col    *mongo.Collection
}

func NewRepoAudit(client mongo.Client) *RepoAudit {
	return &RepoAudit{client: client}
}

func (r *RepoAudit) Coll() *mongo.Collection {
	return r.client.Database("main").Collection("audit")
}

type mongoAudit struct {
	ID           base.UUID     `bson:"_id"`
	RequestType  int           `bson:"request_type"`
	URL          string        `bson:"url"`
	ResponseTime time.Duration `bson:"response_time"`
	CreatedAt    time.Time     `bson:"created_at"`
}

func newFromAudit(audit *audit.Audit) *mongoAudit {
	return &mongoAudit{
		ID:           audit.GetID(),
		RequestType:  audit.GetRequestType(),
		URL:          audit.GetURL(),
		ResponseTime: audit.GetResponseTime(),
		CreatedAt:    time.Now(),
	}
}

func (r *RepoAudit) Create(ctx context.Context, audit *audit.Audit) error {
	_, err := r.col.InsertOne(ctx, newFromAudit(audit))
	return err
}

func (r *RepoAudit) Find(ctx context.Context) ([]*audit.Audit, error) {
	cur, err := r.Coll().Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*audit.Audit
	for cur.Next(ctx) {
		var result mongoAudit
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}

		results = append(results, audit.NewAudit(result.RequestType, result.URL, result.ResponseTime))
	}

	return results, nil
}

func (r *RepoAudit) FindByRequestType(ctx context.Context, requestType int) ([]*audit.Audit, error) {
	cur, err := r.Coll().Find(ctx, bson.D{{Key: "request_type", Value: requestType}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*audit.Audit
	for cur.Next(ctx) {
		var result mongoAudit
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}

		results = append(results, audit.NewAudit(result.RequestType, result.URL, result.ResponseTime))
	}

	return results, nil
}
