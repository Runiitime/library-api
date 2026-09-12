package main

import (
	"encoding/json"
	"fmt"
	"library-api/api"
	"library-api/library"
	"library-api/store"
)

func main() {
	storeData, err := store.LoadData()
	libraryBooks := library.NewLibrary()

	if err == nil {
		var l library.BooksMap
		if err := json.Unmarshal(storeData, &l); err == nil {
			libraryBooks.InitLibrary(l)
		} else {
			panic(err)
		}
	}

	handlers := api.NewHTTPHandlers(libraryBooks)
	server := api.NewServer(handlers)

	if err := server.StartServer(":9091"); err != nil {
		fmt.Println("Failed to start HTTP server:", err)
	}
}
