package db

import (
	"os"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
)

func NewQdrantClient() (*qdrant.Client, error) {
	port, err := strconv.Atoi(os.Getenv("QDRANT_PORT"))
	if err != nil {
		return nil, err
	}
	return qdrant.NewClient(&qdrant.Config{
		Host: os.Getenv("QDRANT_HOST"),
		Port: port,
	})
}
