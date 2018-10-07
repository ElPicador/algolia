package main

import (
	"encoding/csv"
	"io"
)

// TSVReader can be used to read TSV (tab separated values) files.
type TSVReader struct {
	reader   *csv.Reader
	lineChan chan []string
	errChan  chan error
}

// NewTSVReader returns a new TSVReader that begins to read lines from the
// io.Reader passed as parameter. The channels returned by methods Lines and
// Errors must be consumed afterwards.
func NewTSVReader(r io.Reader) *TSVReader {
	ret := &TSVReader{
		reader:   csv.NewReader(r),
		lineChan: make(chan []string),
		errChan:  make(chan error),
	}

	ret.reader.Comma = '\t'
	go ret.Read()
	return ret
}

func (t *TSVReader) Read() {
	defer func() {
		close(t.lineChan)
		close(t.errChan)
	}()

	for {
		record, err := t.reader.Read()
		if err == io.EOF {
			return
		}

		if err != nil {
			t.errChan <- err
		}

		t.lineChan <- record
	}
}

// Lines returns a read-only channel that can be used to fetch lines of a TSV.
func (t *TSVReader) Lines() <-chan []string {
	return t.lineChan
}

// Error returns a read-only channel that can be used to report errors during
// TSV reading.
func (t *TSVReader) Error() <-chan error {
	return t.errChan
}
