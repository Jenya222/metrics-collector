package main

import (
	"net/http"

	"github.com/Jenya222/metrics-collector/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/update/", http.StripPrefix("/update/", handlers.NewUpdateHandler()))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
