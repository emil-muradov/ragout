package aiassitant

import (
	"ai-assistant/internal"
	"context"
	"fmt"
	"os"

	"github.com/sheeiavellie/go-yandexgpt"
)

func main() {
	apiKey := os.Getenv("YANDEXGPT_API_KEY")
	catalogId := os.Getenv("YANDEX_CATALOG_ID")
	client := yandexgpt.New(yandexgpt.CfgApiKey(apiKey))
	vectorDBClient, err := internal.InitVectorDB()
	if err != nil {
		fmt.Println("Error initializing vector DB client")
		return
	}
	defer vectorDBClient.Close()
	internal.InitCollection(context.Background(), vectorDBClient)
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
				Text: "ONE",
			},
		},
	}

	response, err := client.GetCompletion(context.Background(), request)
	if err != nil {
		fmt.Println("Request error")
		return
	}

	fmt.Println(response.Result.Alternatives[0].Message.Text)
}
