package main

import (
	"ai-assistant/internal"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	Status int `json:"status"`
	Message string `json:"message"`
}

type UserRequest struct {
	Msg string `json:"msg"`
}


func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app, err := InitApp(ctx)
	logger := internal.InitLogger()

	if err != nil {
		logger.Error("failed to initialize app", "error", err)
		return
	}

	logger.Info("successfully initialized app")

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "ready to serve!")
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
			err := encoder.Encode(Response{Status: http.StatusInternalServerError, Message: err.Error()})
			if err != nil {
				return
			}
		} else {
			err := encoder.Encode(Response{Status: http.StatusOK, Message: response})
			if err != nil {
				return
			}

		}
	})
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Error("failed to start server", "error", err)
		return
	}
}
