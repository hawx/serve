package main

import (
	"log"
	"net/http"

	"hawx.me/code/serve"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func main() {
	srv := &http.Server{
		Handler: http.HandlerFunc(hello),
	}

	srv.RegisterOnShutdown(func() {
		log.Println("I should really clean up")
	})

	serve.Server("8080", "", srv)
}
