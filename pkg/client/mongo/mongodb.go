package mongodb

import (
	"context"
	"fmt"

	"github.com/Nikitapopov/Habbit/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, logger *logging.Logger, host, port, username, password, database string) (*mongo.Database, error) {
	mongoDBURL := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)

	opts := options.Client().ApplyURI(mongoDBURL)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongoDB %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo db %v", err)
	}

	logger.Info("Pinged your deployment. You successfully connected to MongoDB!")

	return client.Database(database), nil
}
