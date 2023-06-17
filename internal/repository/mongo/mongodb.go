package mongo

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongodb(db string) (*mongo.Database, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("no .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, errors.New("you must set your 'MONGODB_URI' environmental variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(db), nil
}
