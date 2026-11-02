package msg

import "errors"

var ErrBookNotFound = errors.New("book not found")
var ErrJsonWrite = errors.New("failed to write HTTP response")
var ErrWrongQueryParamValue = errors.New("wrong query param value")
var ErrEmptyQuery = errors.New("empty query param")

var ErrQuery = errors.New("query error: ")
var ErrRows = errors.New("rows error: ")
var ErrRowScan = errors.New("scan row error: ")
var ErrExec = errors.New("exec error: ")
var ErrQueryRow = errors.New("query row error: ")
