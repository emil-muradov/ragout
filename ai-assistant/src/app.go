package aiassistant

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"github.com/sheeiavellie/go-yandexgpt"
	"github.com/tmc/langchaingo/schema"
	"golang.org/x/time/rate"
)

type App struct {
	yandexGptClient *yandexgpt.YandexGPTClient
	qdrantClient    *qdrant.Client
	logger          *slog.Logger
}

func InitApp(ctx context.Context) (*App, error) {
	app := &App{}
	app.logger = CreateLogger()
	app.yandexGptClient = yandexgpt.New(yandexgpt.CfgApiKey(os.Getenv("YANDEX_GPT_API_KEY")))
	qdrantClient, err := NewQdrantClient()
	if err != nil {
		return nil, err
	}
	app.qdrantClient = qdrantClient
	err = CreateCollection(ctx, app.qdrantClient)
	if err != nil {
		return nil, err
	}
	err = CreateFieldIndex(ctx, app.qdrantClient)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	uploadKnowledgeBase(ctx, app, &wg)
	return app, nil
}

func (app *App) ProcessUserRequest(ctx context.Context, question string) (string, error) {
	if app == nil {
		return "", errors.New("app is not initialized")
	}
	embedding, err := app.yandexGptClient.GetEmbedding(ctx, yandexgpt.YandexGPTEmbeddingsRequest{
		ModelURI: yandexgpt.MakeEmbModelURI(os.Getenv("YANDEX_CATALOG_ID"), yandexgpt.TextSearchDoc),
		Text:     question,
	})
	if err != nil {
		return "", err
	}
	points, err := Query(ctx, app.qdrantClient, embedding.Embedding)
	if err != nil {
		return "", err
	}
	context := make([]string, len(points))
	for i, point := range points {
		context[i] = point.Payload["description"].GetStringValue()
	}
	templateData := map[string]any{
		"Context": context,
	}
	systemPrompt, err := ParsePrompt("./prompts/system.yml", templateData)
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
				Text: systemPrompt,
			},
			{
				Role: yandexgpt.YandexGPTMessageRoleUser,
				Text: question,
			},
		},
	}
	response, err := app.yandexGptClient.GetCompletion(ctx, request)
	if err != nil {
		return "", err
	}
	return response.Result.Alternatives[0].Message.Text, nil
}

func vectorizeFile(ctx context.Context, app *App, limiter *rate.Limiter, fileContent string) ([]float64, error) {
	err := limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	embedding, err := app.yandexGptClient.GetEmbedding(ctx, yandexgpt.YandexGPTEmbeddingsRequest{
		ModelURI: yandexgpt.MakeEmbModelURI(os.Getenv("YANDEX_CATALOG_ID"), yandexgpt.TextSearchDoc),
		Text:     fileContent,
	})
	if err != nil {
		return nil, err
	}
	return embedding.Embedding, nil
}

func uploadKnowledgeBase(ctx context.Context, app *App, wg *sync.WaitGroup) error {
	limiter := rate.NewLimiter(rate.Every(500*time.Millisecond), 2)
	chunks, err := TextToChunks(ctx, "./knowledge-base/buildings.yml")
	if err != nil {
		return err
	}
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk schema.Document) {
			defer wg.Done()
			vector, err := vectorizeFile(ctx, app, limiter, chunk.PageContent)
			if err != nil {
				return
			}
			err = UploadVectorizedFile(ctx, app.qdrantClient, chunk.PageContent, vector)
			if err != nil {
				return
			}
		}(chunk)
	}
	wg.Wait()
	return nil
}
