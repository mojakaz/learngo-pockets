package gordle_test

import (
	"learngo-pockets/gordle/gordle"
	"testing"
)

func TestReadCorpus(t *testing.T) {
	tt := map[string]struct {
		file   string
		length int
		err    error
	}{
		"English corpus": {
			file:   "../corpus/english.txt",
			length: 156,
			err:    nil,
		},
		"empty corpus": {
			file:   "../corpus/empty.txt",
			length: 0,
			err:    gordle.ErrCorpusIsEmpty,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			words, err := gordle.ReadCorpus(tc.file)
			if len(words) != tc.length {
				t.Errorf("different words length: got %d, want %d", len(words), tc.length)
			}
			if err != tc.err {
				t.Errorf("different error: got %v, want %v", err, tc.err)
			}
		})
	}
}
