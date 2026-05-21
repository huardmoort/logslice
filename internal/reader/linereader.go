package reader

import (
	"bufio"
	"io"
	"os"
)

// LineReader reads lines from a log file with optional offset tracking.
type LineReader struct {
	file    *os.File
	scanner *bufio.Scanner
	lineNum int64
}

// NewLineReader opens a file and returns a LineReader.
func NewLineReader(path string) (*LineReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &LineReader{
		file:    f,
		scanner: scanner,
		lineNum: 0,
	}, nil
}

// NewLineReaderFromReader creates a LineReader from an existing io.Reader.
func NewLineReaderFromReader(r io.Reader) *LineReader {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &LineReader{
		scanner: scanner,
		lineNum: 0,
	}
}

// Next advances to the next line. Returns false when done or on error.
func (lr *LineReader) Next() bool {
	if lr.scanner.Scan() {
		lr.lineNum++
		return true
	}
	return false
}

// Line returns the current line text.
func (lr *LineReader) Line() string {
	return lr.scanner.Text()
}

// LineNumber returns the current 1-based line number.
func (lr *LineReader) LineNumber() int64 {
	return lr.lineNum
}

// Err returns any scanner error (excluding io.EOF).
func (lr *LineReader) Err() error {
	return lr.scanner.Err()
}

// Close closes the underlying file, if any.
func (lr *LineReader) Close() error {
	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}
