package library

import (
	"sync"
)

type BooksMap = map[string]Book

type Library struct {
	Books BooksMap
	mtx   sync.RWMutex
}

func NewLibrary() *Library {
	return &Library{Books: BooksMap{}}
}

func (l *Library) InitLibrary(data BooksMap) {
	if len(data) > 0 {
		l.Books = data
	}
}

func (l *Library) AddBook(b Book, id string) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.Books[id] = b
}

func (l *Library) ChangeCompleted(id string, isCompleted bool) (Book, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	book, ok := l.Books[id]
	if !ok {
		return Book{}, ErrBookNotFound
	}

	book.ChangeCompletedStatus(isCompleted)
	l.Books[id] = book

	return book, nil
}

func (l *Library) FindBook(id string) (Book, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	if book, ok := l.Books[id]; !ok {
		return Book{}, ErrBookNotFound
	} else {
		return book, nil
	}
}

func (l *Library) DeleteBook(id string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if _, ok := l.Books[id]; !ok {
		return ErrBookNotFound
	} else {
		delete(l.Books, id)
		return nil
	}
}

func (l *Library) FilterBooksByAuthor(author string) (BooksMap, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	var books = make(BooksMap)

	for id, book := range l.Books {
		if book.Author == author {
			books[id] = book
		}
	}

	if len(books) == 0 {
		return books, ErrNoBooksFound
	}

	return books, nil
}

func (l *Library) FilterBooksByCompleted(isCompleted bool) (BooksMap, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	var books = make(BooksMap)

	for id, book := range l.Books {
		if book.Completed == isCompleted {
			books[id] = book
		}
	}

	if len(books) == 0 {
		return books, ErrNoBooksFound
	}

	return books, nil
}

func (l *Library) GetAllBooks() (BooksMap, error) {
	if len(l.Books) == 0 {
		return BooksMap{}, ErrLibraryIsEmpty
	}

	var books = make(BooksMap, len(l.Books))

	for id, book := range l.Books {
		books[id] = book
	}

	return books, nil
}

func (l *Library) ValidateBookCreation(title, author string) error {
	bookExist := false

	for _, item := range l.Books {
		if item.Author == author && item.Title == title {
			bookExist = true
			break
		}
	}

	if bookExist == true {
		return ErrBookAlreadyExist
	}

	return nil
}
