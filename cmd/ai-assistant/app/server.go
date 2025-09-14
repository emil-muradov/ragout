package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
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
	router := http.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app, err := InitApp()

	if err != nil {
		log.Fatal("Failed to initialize app")
		return
	}

	router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Ready to serve!")
	})
	router.HandleFunc("POST /ask", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Failed to read body")
			return
		}

		var userRequest UserRequest
		json.Unmarshal(body, &userRequest)
		response, err := app.ProcessUserRequest(ctx, userRequest.Msg)
		encoder := json.NewEncoder(w)

		if err != nil {
			err := encoder.Encode(Response{Status: http.StatusInternalServerError, Message: err.Error()})
			if err != nil {
				log.Fatal("Failed to encode response")
				return
			}
		} else {
			err := encoder.Encode(Response{Status: http.StatusOK, Message: response})
			if err != nil {
				log.Fatal("Failed to encode response")
				return
			}

		}
	})
	http.ListenAndServe(":8080", router)
}
