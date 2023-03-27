package multierr

import (
	"fmt"
)

type FileError struct {
	Filename string
	Line     int
	Column   int
	Err      error
}

func (e *FileError) Unwrap() error {
	return e.Err
}

type Error struct {
	errs []*FileError
}

func (e *Error) Add(filename, in string, loc int, err error) {
	line := 1
	column := 1
	if in != "" && loc != 0 {
		line, column = LineNumber(in, loc)
	}
	e.errs = append(e.errs, &FileError{filename, line, column, err})
}

func (e *Error) Errs() []*FileError {
	return e.errs
}

func (e *Error) Error() string {
	return fmt.Sprintf("multiple errors: %d errors", len(e.errs))
}

func New() *Error {
	return &Error{}
}
