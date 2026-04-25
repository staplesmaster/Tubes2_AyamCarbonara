package main

import (
	"log"
	"net/http"
	"os"

	"github.com/luis/Tubes2_AyamCarbonara/backend/src/handler"
)

func main() {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
