package db

import (
	"context"
	"os"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
)

func ConnectToQdrant() (*qdrant.Client, error) {
	port, err := strconv.Atoi(os.Getenv("QDRANT_PORT"))
	if err != nil {
		return nil, err
	}
	return qdrant.NewClient(&qdrant.Config{
		Host: os.Getenv("QDRANT_HOST"),
		Port: port,
	})
}

func CreateQdrantCollection(ctx context.Context, client *qdrant.Client, collectionName string) {
	client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     256,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}