package gordle

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// corpusError defines a sentinel error.
type corpusError string

// Error is the implementation of the error interface by corpusError.
func (e corpusError) Error() string {
	return string(e)
}

const ErrCorpusIsEmpty = corpusError("corpus is empty")

// ReadCorpus reads the file located at the given path
// and returns a list of words.
func ReadCorpus(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open %q for reading: %w", path, err)
	}

	if len(data) == 0 {
		return nil, ErrCorpusIsEmpty
	}

	// we expect the corpus to be a line- or space-separated list of words
	words := strings.Fields(string(data))

	return words, nil
}

// pickWord return a random word from the corpus.
func pickWord(corpus []string) string {
	index := rand.Intn(len(corpus))
	return corpus[index]
}
