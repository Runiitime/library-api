package handlers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	handlers *HTTPHandlers
}

func NewServer(h *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		handlers: h,
	}
}

func (s *HTTPServer) StartServer(port string) error {
	router := mux.NewRouter()

	router.Path("/books").Methods("POST").HandlerFunc(s.handlers.HandleCreateBook)
	router.Path("/books/{id}").Methods("GET").HandlerFunc(s.handlers.HandleGetBookByID)
	router.Path("/books").Methods("GET").Queries("completed", "{completed}").HandlerFunc(s.handlers.HandleGetUncompletedBooks)
	router.Path("/books").Methods("GET").Queries("author", "{author}").HandlerFunc(s.handlers.HandleGetBooksByAuthor)
	router.Path("/books").Methods("GET").HandlerFunc(s.handlers.HandleGetAllBooks)
	router.Path("/books/{id}").Methods("DELETE").HandlerFunc(s.handlers.HandleDeleteBook)
	router.Path("/books/{id}").Methods("PATCH").HandlerFunc(s.handlers.HandleChangeCompletedStatus)

	if err := http.ListenAndServe(port, router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}
