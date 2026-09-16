package handlers

import (
	"errors"
	"log"
	"net/http"

	_ "library-api/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	)).Methods(http.MethodGet)

	router.Path("/books").Methods("POST").HandlerFunc(s.handlers.HandleCreateBook)
	router.Path("/books/{id}").Methods("GET").HandlerFunc(s.handlers.HandleGetBookByID)
	router.Path("/books").Methods("GET").HandlerFunc(s.handlers.HandleGetAllBooks)
	router.Path("/books/{id}").Methods("DELETE").HandlerFunc(s.handlers.HandleDeleteBook)
	router.Path("/books/{id}").Methods("PATCH").HandlerFunc(s.handlers.HandleChangeCompletedStatus)

	if err := http.ListenAndServe(port, router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Server closed.")
			return nil
		}

		log.Println("Failed to start HTTP server:", err)
		return err
	}

	return nil
}
