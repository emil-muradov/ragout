package main

import (
	"ai-assistant/internal"
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/qdrant/go-client/qdrant"
	"github.com/sheeiavellie/go-yandexgpt"
)

type App struct {
	yandexGptClient *yandexgpt.YandexGPTClient
	vectorDbClient *qdrant.Client
	logger *slog.Logger
}

var apiKey = os.Getenv("YANDEXGPT_API_KEY")
var catalogId = os.Getenv("YANDEX_CATALOG_ID")

func InitApp() (*App, error) {
	app := &App{}
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := godotenv.Load("../.env")

	if err != nil {
		app.logger.Error("failed to load env file", "error", err)
		return nil, err
	}

	app.yandexGptClient = yandexgpt.New(yandexgpt.CfgApiKey(apiKey))
	vectorDBClient, err := internal.InitVectorDB()

	if err != nil {
		app.logger.Error("failed to initialize vector db", "error", err)
		return nil, err
	}

	app.vectorDbClient = vectorDBClient
	defer app.vectorDbClient.Close()
	internal.InitCollection(context.Background(), vectorDBClient, "real_estate")
	app.logger.Info("app initialized")
	return app, nil
}

func (app *App) ProcessUserRequest(ctx context.Context, msg string) (string, error) {
	if app == nil {
		return "", errors.New("app is not initialized")
	}
	request := yandexgpt.YandexGPTRequest{
		ModelURI: yandexgpt.MakeModelURI(catalogId, yandexgpt.YandexGPTLite, yandexgpt.VersionLatest),
		CompletionOptions: yandexgpt.YandexGPTCompletionOptions{
			Stream:      false,
			Temperature: 0.7,
			MaxTokens:   2000,
		},
		Messages: []yandexgpt.YandexGPTMessage{
			{
				Role: yandexgpt.YandexGPTMessageRoleSystem,
				Text: "Ты помощник по недвижимости. Ты помогаешь пользователю найти подходящее жилье.",
			},
			{
				Role: yandexgpt.YandexGPTMessageRoleUser,
				Text: msg,
			},
		},
	}
	response, err := app.yandexGptClient.GetCompletion(ctx, request)

	if err != nil {
		return "", err
	}

	return response.Result.Alternatives[0].Message.Text, nil
}
