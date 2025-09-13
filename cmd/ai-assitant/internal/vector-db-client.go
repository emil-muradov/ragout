package internal

import (
	"context"
	"os"

	"github.com/qdrant/go-client/qdrant"
)

func InitVectorDB() (*qdrant.Client, error) {
	return qdrant.NewClient(&qdrant.Config{
		Host: os.Getenv("QDRANT_HOST"),
		APIKey: os.Getenv("QDRANT_API_KEY"),
	})
}

func InitCollection(ctx context.Context, client *qdrant.Client) {
	client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "real_estate",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     4,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}