package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"library-api/helpers"
	"library-api/library"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	library *library.Library
}

func NewHTTPHandlers(l *library.Library) *HTTPHandlers {
	return &HTTPHandlers{
		library: l,
	}
}

/*
pattern: /books
method: POST
info: JSON in HTTP request body

succeed:

  - status code: 201 Created
  - response body: JSON represent created book

failed:

  - status code: 400, 409, 500
  - response body: JSON with error
*/
func (h *HTTPHandlers) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	var data BookDTO

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		DoError(library.ErrJSONInmarshal, w, http.StatusInternalServerError)
		return
	}

	if err := h.library.ValidateBookCreation(data.Title, data.Author); err != nil {
		DoError(err, w)
		return
	}

	newBook := library.NewBook(data.Title, data.Author, data.Pages)

	h.library.AddBook(*newBook, helpers.GenerateID())
	b, err := json.MarshalIndent(newBook, "", "    ")

	if err != nil {
		DoError(library.ErrJSONInmarshal, w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

/*
pattern: /books/{id}
method: GET
info: pattern

succeed:

  - status code: 200 Ok
  - response body: JSON represent founded book

failed:

  - status code: 400, 404, 500
  - response body: JSON with error
*/
func (h *HTTPHandlers) HandleGetBook(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	book, err := h.library.FindBook(id)

	if err != nil {
		if errors.Is(err, library.ErrBookNotFound) {
			DoError(err, w, http.StatusNotFound)
			return
		}
		DoError(err, w, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(book, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		DoError(err, w, http.StatusInternalServerError)
		return
	}
}

/*
pattern: /books/{id}
method: DELETE
info: pattern

succeed:

  - status code: 204 No content
  - response body: -

failed:

  - status code: 400, 404, 500
  - response body: JSON with error
*/
func (h *HTTPHandlers) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]

	if err := h.library.DeleteBook(id); err != nil {
		DoError(err, w, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
pattern: /books
method: GET
info: -

succeed:

  - status code: 200 Ok
  - response body: - JSON represent found books

failed:

  - status code: 400, 500
  - response body: JSON with error
*/
func (h *HTTPHandlers) HandleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.library.GetAllBooks()

	if err != nil {
		if !errors.Is(err, library.ErrLibraryIsEmpty) {
			DoError(err, w, http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(books, "", "    ")

	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		DoError(err, w, http.StatusInternalServerError)
		return
	}
}

/*
pattern: /books/{id}
method: PATCH
info: pattern + JSON in request body { is_completed: bool }

succeed:
  - status code: 200 Ok
  - response body: JSON represent changed book

failed:
  - status code: 400, 404, 409, 500
  - response body: JSON with error
*/
func (h *HTTPHandlers) HandleChangeCompletedStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	decoder := json.NewDecoder(r.Body)

	var status BookStatusDTO

	if err := decoder.Decode(&status); err != nil {
		DoError(library.ErrJSONInmarshal, w, http.StatusInternalServerError)
		return
	}

	book, err := h.library.ChangeCompleted(id, status.Completed)

	if err != nil {
		DoError(err, w, http.StatusNotFound)
		return
	}

	b, err := json.MarshalIndent(book, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		DoError(err, w, http.StatusInternalServerError)
		return
	}
}

/*
pattern: /books?completed=false
method: GET
info: query params

succeed:
- status code: 200 Ok
- response body: JSON represent completed books

failed:
- status code: 400, 404, 409, 500
- response body: JSON with error
*/
func (h *HTTPHandlers) HandleGetUncompletedBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	param := query.Get("completed")

	if param == "" {
		DoError(library.ErrEmptyQuery, w)
		return
	}

	completed, err := strconv.ParseBool(param)

	if err != nil {
		DoError(library.ErrWrongQueryParamValue, w)
		return
	}

	books, err := h.library.FilterBooksByCompleted(completed)

	if err != nil {
		if errors.Is(library.ErrNoBooksFound, err) {
			DoError(err, w, http.StatusOK)
			return
		}
		panic(err)
	}

	data, err := json.MarshalIndent(books, "", "    ")

	if err != nil {
		panic(err)
	}

	if _, err := w.Write(data); err != nil {
		fmt.Println(library.ErrJsonWrite, err)
		return
	}
}

/*
pattern: /books?completed=false
method: GET
info: query params

succeed:
- status code: 200 Ok
- response body: JSON represent completed books

failed:
- status code: 400, 404, 409, 500
- response body: JSON with error
*/
func (h *HTTPHandlers) HandleGetBooksByAuthor(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	param := query.Get("author")

	if param == "" {
		DoError(library.ErrEmptyQuery, w)
		return
	}

	books, err := h.library.FilterBooksByAuthor(param)

	if err != nil {
		if errors.Is(library.ErrNoBooksFound, err) {
			DoError(err, w, http.StatusOK)
			return
		}
		panic(err)
	}

	data, err := json.MarshalIndent(books, "", "    ")

	if err != nil {
		panic(err)
	}

	if _, err := w.Write(data); err != nil {
		fmt.Println(library.ErrJsonWrite, err)
		return
	}
}
