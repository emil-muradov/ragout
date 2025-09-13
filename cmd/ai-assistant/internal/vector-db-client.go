package internal

import (
	"context"
	"os"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
)

func InitVectorDB() (*qdrant.Client, error) {
	port, err := strconv.Atoi(os.Getenv("QDRANT_PORT"))
	if err != nil {
		return nil, err
	}
	return qdrant.NewClient(&qdrant.Config{
		Host: os.Getenv("QDRANT_HOST"),
		Port: port,
	})
}

func InitCollection(ctx context.Context, client *qdrant.Client, collectionName string) {
	client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     4,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}