package handlers

import (
	"encoding/json"
	"errors"
	"library-api/models"
	"library-api/storage"
	bookMSG "library-api/storage/msg"
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

// HandleCreateBook Create book
//
//	@Summary      Create a book
//	@Description  create a book
//	@Tags		  book
//	@Accept       json
//	@Produce      json
//	@Param		  book	body	 models.Book	true	"book info"
//	@Success      201  {object}  models.Book
//	@Failure      500  {object}  models.ErrorDTO
//	@Router       /books [post]
func (h *HTTPHandlers) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	var data models.Book

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}

	id, err := h.queries.CreateBook(r.Context(), data)
	if err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}

	data.ID = id

	log.Println(bookMSG.BookWasCreated)

	b, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}
}

// HandleDeleteBook Delete book
//
//	@Summary      Delete a book
//	@Description  delete a book by id
//	@Tags		  book
//	@Param		  id	path	 int	true	"book ID"	Format(int)
//	@Success      204  {}  No content
//	@Failure      404  {object}  models.ErrorDTO
//	@Failure      500  {object}  models.ErrorDTO
//	@Router       /books/{id}    [delete]
func (h *HTTPHandlers) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	bookID, err := strconv.Atoi(id)
	if err != nil {
		models.DoError(err, w, http.StatusNotFound)
		return
	}

	if err := h.queries.DeleteBook(r.Context(), bookID); err != nil {
		if errors.Is(err, bookMSG.ErrBookNotFound) {
			models.DoError(bookMSG.ErrBookNotFound, w, http.StatusNotFound)
			return
		}
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetBookByID Get book by id
//
//		@Summary      Get a book by id
//		@Description  get a book by id
//		@Tags		  book
//		@Produce      json
//	 	@Param 		  id path int true "book ID"	Format(int)
//		@Success      200  {object}  models.Book
//		@Failure      404  {object}  models.ErrorDTO
//		@Failure      500  {object}  models.ErrorDTO
//		@Router       /books/{id} [get]
func (h *HTTPHandlers) HandleGetBookByID(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]

	n, _ := strconv.Atoi(id)
	books, err := h.queries.SelectBooksByID(r.Context(), []int{n})

	if err != nil {
		models.DoError(err, w)
		return
	}

	if len(books) == 0 {
		models.DoError(bookMSG.ErrBookNotFound, w, http.StatusNotFound)
		return
	}

	b, err := json.MarshalIndent(books[0], "", "    ")
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := w.Write(b); err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}
}

// HandleChangeCompletedStatus Change book status "completed"
//
//	@Summary      Change a book status "completed"
//	@Description  change a book status "completed"
//	@Tags		  book
//	@Accept       json
//	@Param 		  id path int true "Book ID"	Format(int)
//	@Param 		  completed body  models.BookStatusDTO true "completed value"
//	@Success      204  {}       No content
//	@Failure      404  {object}  models.ErrorDTO
//	@Failure      500  {object}  models.ErrorDTO
//	@Router       /books/{id} [patch]
func (h *HTTPHandlers) HandleChangeCompletedStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	decoder := json.NewDecoder(r.Body)

	var status models.BookStatusDTO

	if err := decoder.Decode(&status); err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}

	n, err := strconv.Atoi(id)

	if err != nil {
		models.DoError(err, w)
		return
	}

	if err := h.queries.UpdateBookStatus(r.Context(), n, status.Completed); err != nil {
		if errors.Is(err, bookMSG.ErrBookNotFound) {
			models.DoError(bookMSG.ErrBookNotFound, w, http.StatusNotFound)
			return
		}
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetAllBooks Get all books
//
//	@Summary      Get all books
//	@Description  get all books, optionally filtered by status or author
//	@Tags		  books
//	@Produce      json
//	@Param 		  completed query string    false "string valid"
//	@Param		  author query	  string	false	"Search books by author"
//	@Success      200  {array}   []models.Book
//	@Failure      500  {object}  models.ErrorDTO
func (h *HTTPHandlers) HandleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	completedParam := query.Get("completed")
	authorParam := query.Get("author")

	var books []models.Book
	var err error

	if len(query) == 0 {
		books, err = h.queries.SelectAllBooks(r.Context())
	}

	if len(query) == 1 {
		if completedParam != "" {

			_, err := strconv.ParseBool(completedParam)

			if err != nil {
				models.DoError(bookMSG.ErrWrongQueryParamValue, w)
				return
			}

			books, err = h.queries.SelectBooksByParams(r.Context(), "completed", completedParam)
		}

		if authorParam != "" {
			books, err = h.queries.SelectBooksByParams(r.Context(), "author", authorParam)
		}
	}

	if err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(books, "", "    ")

	if err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		models.DoError(err, w, http.StatusInternalServerError)
		return
	}
}
