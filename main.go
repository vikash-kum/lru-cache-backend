package main

import (
	"log"
	"lru-cache-gin/routes"
)

func main() {
	// Initialize the Gin router
	r := routes.SetupRouter()

	// Start the server
	log.Println("Starting server on :8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
