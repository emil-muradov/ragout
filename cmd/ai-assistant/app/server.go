package main

import (
	"context"
	"io"
	"log"
	"net/http"
)


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
		app.ProcessUserRequest(ctx, string(body))
	})
	http.ListenAndServe(":8080", router)
}
