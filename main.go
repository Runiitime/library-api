package main

import (
	"context"
	"library-api/config"
	"library-api/handlers"
	"library-api/logger"
	"library-api/storage"
	"log"
)

// @title           Library API
// @version         1.0
// @description     API Server for Library application
// @termsOfService  http://swagger.io/terms/

// @host      localhost:9091
// @BasePath  /
func main() {
	logger.LoadLogger()

	ctx := context.Context(context.Background())
	cfg := config.GetConfig()

	conn, err := storage.CreateConnection(ctx, cfg)

	bookQuery := storage.NewBookQuery("library", conn)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to database")

	hdls := handlers.NewHTTPHandlers(bookQuery)
	server := handlers.NewServer(hdls)

	log.Println("Listening on port 9091")
	if err := server.StartServer(":9091"); err != nil {
		log.Fatalln("Failed to start HTTP server:", err)
	}
}
