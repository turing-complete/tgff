// Package tgff provides a parser for the TGFF (Task Graphs For Free) format,
// which is a format for storing task graphs and accompanying data used in
// scheduling and allocation research.
//
// http://ziyang.eecs.umich.edu/~dickrp/tgff
package tgff

import (
	"io"
	"os"
)

// Parse reads the content of a TGFF file (*.tgff), generated by the tgff
// command-line tool from a TGFFOPT file (*.tgffopt), and returns its
// representation in a Result struct.
func Parse(reader io.Reader) (Result, error) {
	abort := make(chan bool, 1)

	lexer, stream := newLexer(reader, abort)
	parser, success, failure := newParser(stream, abort)

	go lexer.run()
	go parser.run()

	select {
	case result := <-success:
		return result, nil
	case err := <-failure:
		return Result{}, err
	}
}

// ParseFile works exactly as Parse but takes a path to a TGFF file instead of
// an io.Reader.
func ParseFile(path string) (Result, error) {
	file, err := os.Open(path)

	if err != nil {
		return Result{}, err
	}

	defer file.Close()

	return Parse(file)
}
