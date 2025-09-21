package app

import (
	"context"
	"errors"
	"os"
	"time"

	"ragout/cmd/ai-assistant/internal/db"
	log "ragout/cmd/ai-assistant/internal/logger"
	"ragout/cmd/ai-assistant/internal/utils"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
	"github.com/sheeiavellie/go-yandexgpt"
	"github.com/tmc/langchaingo/schema"
	"golang.org/x/time/rate"
)

const QDRANT_COLLECTION_NAME = "real_estate"

type App struct {
	yandexGptClient *yandexgpt.YandexGPTClient
	qdrantClient    *qdrant.Client
}

func InitApp(ctx context.Context) (*App, error) {
	app := &App{}
	app.yandexGptClient = yandexgpt.New(yandexgpt.CfgApiKey(os.Getenv("YANDEXGPT_API_KEY")))
	qdrantClient, err := db.NewQdrantClient()
	if err != nil {
		return nil, err
	}
	app.qdrantClient = qdrantClient
	exists, err := app.qdrantClient.CollectionExists(ctx, QDRANT_COLLECTION_NAME)
	if err != nil {
		return nil, err
	}
	if exists {
		return app, nil
	}
	app.qdrantClient.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: QDRANT_COLLECTION_NAME,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     256,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	app.qdrantClient.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
		CollectionName: QDRANT_COLLECTION_NAME,
		FieldName:      "description",
		FieldType:      qdrant.FieldType_FieldTypeText.Enum(),
	})
	client := app.yandexGptClient
	chunks, err := utils.TextToChunks(ctx, "./internal/prompts/buildings.txt")
	if err != nil {
		return nil, err
	}
	limiter := rate.NewLimiter(rate.Every(500*time.Millisecond), 2)
	for _, chunk := range chunks {
		go func(doc schema.Document) {
			err := limiter.Wait(ctx)
			if err != nil {
				log.Error("rate limiter: failed to wait", "error", err)
				return
			}
			embedding, err := client.GetEmbedding(ctx, yandexgpt.YandexGPTEmbeddingsRequest{
				ModelURI: yandexgpt.MakeEmbModelURI(os.Getenv("YANDEX_CATALOG_ID"), yandexgpt.TextSearchDoc),
				Text:     doc.PageContent,
			})
			if err != nil {
				log.Error("failed to get embedding", "error", err)
				return
			}
			_, err = qdrantClient.Upsert(ctx, &qdrant.UpsertPoints{
				CollectionName: QDRANT_COLLECTION_NAME,
				Points: []*qdrant.PointStruct{
					{
						Id:      qdrant.NewIDUUID(uuid.NewString()),
						Vectors: qdrant.NewVectors(utils.ConvertFloat64VectorToFloat32Vector(embedding.Embedding)...),
						Payload: qdrant.NewValueMap(map[string]any{
							"description": doc.PageContent,
						}),
					},
				},
			})
			if err != nil {
				log.Error("failed to upsert chunk", "error", err)
				return
			}
		}(chunk)
	}
	return app, nil
}

func (app *App) ProcessUserRequest(ctx context.Context, msg string) (string, error) {
	if app == nil {
		return "", errors.New("app is not initialized")
	}
	embedding, err := app.yandexGptClient.GetEmbedding(ctx, yandexgpt.YandexGPTEmbeddingsRequest{
		ModelURI: yandexgpt.MakeEmbModelURI(os.Getenv("YANDEX_CATALOG_ID"), yandexgpt.TextSearchDoc),
		Text:     msg,
	})
	if err != nil {
		return "", err
	}
	searchResult, err := app.qdrantClient.Query(ctx, &qdrant.QueryPoints{
		CollectionName: QDRANT_COLLECTION_NAME,
		Query:          qdrant.NewQuery(utils.ConvertFloat64VectorToFloat32Vector(embedding.Embedding)...),
	})
	if err != nil {
		return "", err
	}
	pointsIds := make([]*qdrant.PointId, len(searchResult))
	for i, result := range searchResult {
		pointsIds[i] = qdrant.NewIDUUID(result.Id.GetUuid())
	}
	points, err := app.qdrantClient.Get(ctx, &qdrant.GetPoints{
		CollectionName: QDRANT_COLLECTION_NAME,
		Ids:            pointsIds,
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return "", err
	}
	context := make([]string, len(points))
	for i, point := range points {
		context[i] = point.Payload["description"].GetStringValue()
	}
	templateData := map[string]any{
		"Context":  context,
		"Question": msg,
	}
	augmentedPrompt, err := utils.ParsePrompt("./internal/prompts/system.txt", templateData)
	if err != nil {
		return "", err
	}
	request := yandexgpt.YandexGPTRequest{
		ModelURI: yandexgpt.MakeModelURI(os.Getenv("YANDEX_CATALOG_ID"), yandexgpt.YandexGPTLite, yandexgpt.VersionLatest),
		CompletionOptions: yandexgpt.YandexGPTCompletionOptions{
			Stream:      false,
			Temperature: 0.7,
			MaxTokens:   2000,
		},
		Messages: []yandexgpt.YandexGPTMessage{
			{
				Role: yandexgpt.YandexGPTMessageRoleSystem,
				Text: augmentedPrompt,
			},
		},
	}
	response, err := app.yandexGptClient.GetCompletion(ctx, request)
	if err != nil {
		return "", err
	}
	return response.Result.Alternatives[0].Message.Text, nil
}
