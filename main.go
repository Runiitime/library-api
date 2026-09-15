package main

import (
	"context"
	"library-api/config"
	"library-api/handlers"
	"library-api/storage"
	"log"
)

func main() {
	ctx := context.Context(context.Background())
	cfg := config.GetConfig()

	conn, err := storage.CreateConnection(ctx, cfg)

	bookQuery := storage.NewBookQuery("library", conn)
	if err != nil {
		panic(err)
	}

	log.Println("Connected to database")

	hdls := handlers.NewHTTPHandlers(bookQuery)
	server := handlers.NewServer(hdls)

	if err := server.StartServer(":9091"); err != nil {
		log.Fatalln("Failed to start HTTP server:", err)
	}
}
