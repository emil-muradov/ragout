package aiassistant

import (
	"context"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

const QDRANT_COLLECTION_NAME = "real_estate"

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

func CreateCollection(ctx context.Context, client *qdrant.Client) error {
	exists, err := client.CollectionExists(ctx, QDRANT_COLLECTION_NAME)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: QDRANT_COLLECTION_NAME,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     256,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

func CreateFieldIndex(ctx context.Context, client *qdrant.Client) error {
	_, err := client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
		CollectionName: QDRANT_COLLECTION_NAME,
		FieldName:      "description",
		FieldType:      qdrant.FieldType_FieldTypeText.Enum(),
	})
	return err
}

func UploadVectorizedFile(ctx context.Context, client *qdrant.Client, fileContent string, vector []float64) error {
	_, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: QDRANT_COLLECTION_NAME,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDUUID(uuid.NewString()),
				Vectors: qdrant.NewVectors(ConvertFloat64VectorToFloat32Vector(vector)...),
				Payload: qdrant.NewValueMap(map[string]any{
					"description": fileContent,
				}),
			},
		},
	})
	return err
}

func Query(ctx context.Context, client *qdrant.Client, input []float64) ([]*qdrant.RetrievedPoint, error) {
	searchResult, err := client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: QDRANT_COLLECTION_NAME,
		Query:          qdrant.NewQuery(ConvertFloat64VectorToFloat32Vector(input)...),
	})
	if err != nil {
		return nil, err
	}
	pointsIds := make([]*qdrant.PointId, len(searchResult))
	for i, result := range searchResult {
		pointsIds[i] = qdrant.NewIDUUID(result.Id.GetUuid())
	}
	points, err := client.Get(ctx, &qdrant.GetPoints{
		CollectionName: QDRANT_COLLECTION_NAME,
		Ids:            pointsIds,
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, err
	}
	return points, nil
}
