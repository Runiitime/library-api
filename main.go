package main

import (
	"context"
	"fmt"
	"library-api/config"
	"library-api/storage"
)

func main() {
	//libraryBooks := library.NewLibrary()

	//handlers := api.NewHTTPHandlers(libraryBooks)
	//server := api.NewServer(handlers)
	//
	//if err := server.StartServer(":9091"); err != nil {
	//	fmt.Println("Failed to start HTTP server:", err)
	//}
	ctx := context.Context(context.Background())
	cfg := config.GetConfig()

	conn, err := storage.CreateConnection(ctx, cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database")

	if err := storage.CreateTable(conn, ctx, "library"); err != nil {
		panic(err)
	}

	fmt.Println("Created table library")
}
