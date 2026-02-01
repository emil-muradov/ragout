package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	aiassistant "github.com/emil-muradov/ragout/ai-assistant/src"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Request struct {
	Question string `json:"question"`
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app, err := aiassistant.InitApp(ctx)
	if err != nil {
		log.Println("failed to initialize app", "error", err)
		return
	}
	log.Println("app initialized")
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "OK")
	})
	router.HandleFunc("POST /ask", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		var request Request
		encoder := json.NewEncoder(w)
		err = json.Unmarshal(body, &request)
		if err != nil {
			encoder.Encode(Response{Status: http.StatusBadRequest, Message: err.Error()})
			return
		}
		response, err := app.ProcessUserRequest(ctx, request.Question)
		if err != nil {
			encoder.Encode(Response{Status: http.StatusInternalServerError, Message: err.Error()})
		} else {
			encoder.Encode(Response{Status: http.StatusOK, Message: response})
		}
	})
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Println("failed to start server", "error", err)
		return
	}
}
