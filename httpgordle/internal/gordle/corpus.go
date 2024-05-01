package gordle

import (
	"math/rand"
	"strings"
)

// corpusError defines a sentinel error.
type corpusError string

// Error is the implementation of the error interface by corpusError.
func (e corpusError) Error() string {
	return string(e)
}

const ErrCorpusIsEmpty = corpusError("corpus is empty")

// ReadCorpus reads the given corpus and returns a list of words.
func ReadCorpus(data string) ([]string, error) {
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
