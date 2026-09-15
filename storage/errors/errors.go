package errors

import "errors"

// Book not found
var ErrBookNotFound = errors.New("Book not found")

// Book already exist
var ErrBookAlreadyExist = errors.New("Book already exist")

// No books were found
var ErrNoBooksFound = errors.New("No books were found")

// Library is empty"
var ErrLibraryIsEmpty = errors.New("Library is empty")

// Failed to write HTTP response
var ErrJsonWrite = errors.New("Failed to write HTTP response")

// Wrong query param value
var ErrWrongQueryParamValue = errors.New("Wrong query param value")

// Empty query param
var ErrEmptyQuery = errors.New("Empty query param")

// Empty query param
var ErrJSONUnmarshal = errors.New("Failed to json unmarshal")
