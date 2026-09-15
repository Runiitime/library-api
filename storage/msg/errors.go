package msg

import "errors"

var ErrBookNotFound = errors.New("Book not found")
var ErrJsonWrite = errors.New("Failed to write HTTP response")
var ErrWrongQueryParamValue = errors.New("Wrong query param value")
var ErrEmptyQuery = errors.New("Empty query param")

var ErrQuery = errors.New("Query error:")
var ErrRows = errors.New("Rows error:")
var ErrRowScan = errors.New("Scan row error:")
var ErrExec = errors.New("Exec error:")
var ErrQueryRow = errors.New("Query row error:")
