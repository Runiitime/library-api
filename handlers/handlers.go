package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"library-api/models"
	"library-api/storage"
	bookErr "library-api/storage/errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	queries *storage.BookQuery
}

func NewHTTPHandlers(queries *storage.BookQuery) *HTTPHandlers {
	return &HTTPHandlers{
		queries: queries,
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
	var data models.Book

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		models.DoError(bookErr.ErrJSONUnmarshal, w, http.StatusInternalServerError)
		return
	}

	id, err := h.queries.CreateBook(r.Context(), data)
	if err != nil {
		log.Panic(err)
	}

	data.ID = id

	log.Println("Книга была создана")

	b, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		models.DoError(bookErr.ErrJSONUnmarshal, w, http.StatusInternalServerError)
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
func (h *HTTPHandlers) HandleGetBookByID(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]

	n, _ := strconv.Atoi(id)
	books, err := h.queries.SelectBooksByID(r.Context(), []int{n})

	if err != nil {
		models.DoError(err, w)
		return
	}

	if len(books) == 0 {
		models.DoError(bookErr.ErrBookNotFound, w, http.StatusNotFound)
		return
	}

	b, err := json.MarshalIndent(books[0], "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
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
	bookID, err := strconv.Atoi(id)
	if err != nil {
		models.DoError(err, w, http.StatusNotFound)
		return
	}

	if err := h.queries.DeleteBook(r.Context(), bookID); err != nil {
		if errors.Is(err, bookErr.ErrBookNotFound) {
			models.DoError(bookErr.ErrBookNotFound, w, http.StatusNotFound)
			return
		}
		models.DoError(err, w, http.StatusInternalServerError)
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
	books, err := h.queries.SelectAllBooks(r.Context())
	if err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}

	if len(books) == 0 {
		models.DoError(bookErr.ErrLibraryIsEmpty, w, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(books, "", "    ")

	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}
}

/*
pattern: /books/{id}
method: PATCH
info: pattern + JSON in request body { is_completed: bool }

succeed:
  - status code: 204 No content
  - response body: -

failed:
  - status code: 400, 404, 409, 500
  - response body: JSON with error
*/
func (h *HTTPHandlers) HandleChangeCompletedStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	decoder := json.NewDecoder(r.Body)

	var status models.BookStatusDTO

	if err := decoder.Decode(&status); err != nil {
		models.DoError(bookErr.ErrJSONUnmarshal, w, http.StatusInternalServerError)
		return
	}

	n, err := strconv.Atoi(id)

	if err != nil {
		models.DoError(err, w, http.StatusBadRequest)
	}

	if err := h.queries.UpdateBookStatus(r.Context(), n, status.Completed); err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		models.DoError(bookErr.ErrEmptyQuery, w)
		return
	}

	_, err := strconv.ParseBool(param)

	if err != nil {
		models.DoError(bookErr.ErrWrongQueryParamValue, w)
		return
	}

	books, err := h.queries.SelectBooksByParams(r.Context(), "completed", param)
	if err != nil {
		if errors.Is(bookErr.ErrNoBooksFound, err) {
			models.DoError(err, w, http.StatusOK)
			return
		}
		panic(err)
	}

	data, err := json.MarshalIndent(books, "", "    ")

	if err != nil {
		panic(err)
	}

	if _, err := w.Write(data); err != nil {
		fmt.Println(bookErr.ErrJsonWrite, err)
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
		models.DoError(bookErr.ErrEmptyQuery, w)
		return
	}

	books, err := h.queries.SelectBooksByParams(r.Context(), "author", param)
	if err != nil {
		if errors.Is(bookErr.ErrNoBooksFound, err) {
			models.DoError(err, w, http.StatusOK)
			return
		}
		panic(err)
	}

	data, err := json.MarshalIndent(books, "", "    ")

	if err != nil {
		panic(err)
	}

	if _, err := w.Write(data); err != nil {
		fmt.Println(bookErr.ErrJsonWrite, err)
		return
	}
}
