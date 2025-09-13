package main

import (
	"ai-assistant/internal"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/qdrant/go-client/qdrant"
	"github.com/sheeiavellie/go-yandexgpt"
)

type App struct {
	yandexGptClient *yandexgpt.YandexGPTClient
	vectorDbClient *qdrant.Client
}

var apiKey = os.Getenv("YANDEXGPT_API_KEY")
var catalogId = os.Getenv("YANDEX_CATALOG_ID")

func InitApp() (*App, error) {
	app := &App{}
	app.yandexGptClient = yandexgpt.New(yandexgpt.CfgApiKey(apiKey))
	vectorDBClient, err := internal.InitVectorDB()

	if err != nil {
		log.Fatal("Failed to initialize vector db")
		return nil, err
	}

	app.vectorDbClient = vectorDBClient
	defer app.vectorDbClient.Close()
	internal.InitCollection(context.Background(), vectorDBClient, "real_estate")
	return app, nil
}

func (app *App) ProcessUserRequest(ctx context.Context, input string) (string, error) {
	if app == nil {
		return "", errors.New("App is not initialized")
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
				Text: input,
			},
		},
	}
	response, err := app.yandexGptClient.GetCompletion(ctx, request)

	if err != nil {
		return "", errors.New("Request error")
	}

	fmt.Println(response.Result.Alternatives[0].Message.Text)
	return response.Result.Alternatives[0].Message.Text, nil
}
