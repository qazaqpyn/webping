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
	col *mongo.Collection
}

func NewRepoAudit(db *mongo.Database, coll string) *RepoAudit {
	return &RepoAudit{
		col: db.Collection(coll),
	}
}

type mongoAudit struct {
	ID           base.UUID     `bson:"_id"`
	RequestType  int           `bson:"request_type"`
	URL          string        `bson:"url"`
	ResponseTime time.Duration `bson:"response_time"`
	CreatedAt    time.Time     `bson:"created_at"`
}

type mongoAuditGroup struct {
	ID    int `bson:"_id"`
	Total int `bson:"total"`
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

func (r *RepoAudit) Find(ctx context.Context) ([]*audit.MongoAuditGroup, error) {
	// agregate by request type and get total number of requests
	cur, err := r.col.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "request_type", Value: "$request_type"},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []*audit.MongoAuditGroup
	for cur.Next(ctx) {
		var result mongoAuditGroup
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}

		results = append(results, &audit.MongoAuditGroup{
			ID:    result.ID,
			Total: result.Total,
		})
	}

	return results, nil
}

func (r *RepoAudit) FindByRequestType(ctx context.Context, requestType int) ([]*audit.Audit, error) {
	cur, err := r.col.Find(ctx, bson.D{{Key: "request_type", Value: requestType}})
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
