package main

import (
	log "ai-assistant/internal/logger"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type UserRequest struct {
	Msg string `json:"msg"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := log.CreateLogger()
	app, err := InitApp(ctx)
	if err != nil {
		logger.Error("failed to initialize app", "error", err)
		return
	}
	logger.Info("app initialized")
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "ready to serve")
	})
	router.HandleFunc("POST /ask", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		var userRequest UserRequest
		json.Unmarshal(body, &userRequest)
		response, err := app.ProcessUserRequest(ctx, userRequest.Msg)
		encoder := json.NewEncoder(w)
		if err != nil {
			encoder.Encode(Response{Status: http.StatusInternalServerError, Message: err.Error()})
		} else {
			encoder.Encode(Response{Status: http.StatusOK, Message: response})
		}
	})
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Error("failed to start server", "error", err)
		return
	}
}
